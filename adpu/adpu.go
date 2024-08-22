package adpu

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ntchjb/ledger-go/device"
	"github.com/ntchjb/ledger-go/log"
)

const (
	// ADPU header length for a block
	ADPU_BLOCK_HEADER_LENGTH uint16 = 5
	FRAME_HEADER_TAG         uint8  = 0x05
)

var (
	ErrSWNotOK = errors.New("SW not OK")
)

type Response struct {
	Length uint16
	Data   []byte
}

type Protocol interface {
	Exchange(ctx context.Context, command []byte) ([]byte, error)
	Send(ctx context.Context, cla, ins, p1, p2 uint8, data []byte) ([]byte, uint16, error)
}

type protocolImpl struct {
	Device     device.Device
	logger     *slog.Logger
	channel    uint16
	packetSize uint16

	exchangeLock sync.RWMutex
}

func NewProtocol(device device.Device, channel uint16, logger *slog.Logger) Protocol {
	return &protocolImpl{
		Device:     device,
		logger:     logger,
		channel:    channel,
		packetSize: 64,
	}
}

func (a *protocolImpl) createDataFrames(data []byte) [][]byte {
	a.logger.Debug("Creating data frames", "length", len(data))
	var blocks [][]byte
	dataLength := uint16(len(data))
	blockSizeWithoutHeader := a.packetSize - ADPU_BLOCK_HEADER_LENGTH
	numBlocks := dataLength / blockSizeWithoutHeader
	if dataLength%blockSizeWithoutHeader != 0 {
		numBlocks++
	}
	a.logger.Debug("Data frame properties", "numBlocks", numBlocks, "blockSizeWithoutHeader", blockSizeWithoutHeader)

	padding := make([]byte, numBlocks*blockSizeWithoutHeader-dataLength+1)
	dataLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(dataLengthBytes, dataLength)

	// Final data should be
	// [dataLength..., data..., padding...]
	data = append(dataLengthBytes, data...)
	data = append(data, padding...)

	a.logger.Debug("Final data before convert to frames", "data", log.HexDisplay(data))

	for i := uint16(0); i < numBlocks; i++ {
		header := make([]byte, 5)
		binary.BigEndian.PutUint16(header[:2], a.channel)
		header[2] = FRAME_HEADER_TAG
		binary.BigEndian.PutUint16(header[3:5], i)

		blocks = append(blocks, append(header, data[i*blockSizeWithoutHeader:(i+1)*blockSizeWithoutHeader]...))
	}

	a.logger.Debug("Frames", "blockCount", len(blocks))

	return blocks
}

func (a *protocolImpl) reduceDataFrames(res Response, expectedSequence uint16, block []byte) (Response, error) {
	a.logger.Debug("Reducing a frame", "blockLength", len(block), "expectedSequence", expectedSequence, "resLength", res.Length)
	if len(block) < 5 {
		return res, fmt.Errorf("chunk data is too small, expected >=5, but got %d", len(block))
	}
	channel := binary.BigEndian.Uint16(block[:2])
	tag := block[2]
	sequence := binary.BigEndian.Uint16(block[3:5])
	if channel != a.channel {
		return res, fmt.Errorf("channel does not match, expected %d, got %d", a.channel, channel)
	}
	if tag != byte(FRAME_HEADER_TAG) {
		return res, fmt.Errorf("tag does not match, expected %d, got %d", FRAME_HEADER_TAG, tag)
	}
	if sequence != expectedSequence {
		return res, fmt.Errorf("sequence does not match, expected %d, got %d", expectedSequence, sequence)
	}

	// Read length for the first block
	if len(res.Data) == 0 {
		res.Length = binary.BigEndian.Uint16(block[5:7])
		res.Data = append(res.Data, block[7:]...)
	} else {
		res.Data = append(res.Data, block[5:]...)
	}

	// Remove padding from res.Data at the last block
	if len(res.Data) > int(res.Length) {
		res.Data = res.Data[:res.Length]
	}

	return res, nil
}

// Send ADPU command to Device via Ledger's HID report scheme
// Given command
func (a *protocolImpl) Exchange(ctx context.Context, command []byte) ([]byte, error) {
	a.exchangeLock.Lock()
	defer a.exchangeLock.Unlock()

	a.logger.Debug("ADPU Command", "command", log.HexDisplay(command))

	// #1: Send ADPU command to device, in blocks
	blocks := a.createDataFrames(command)
	for _, block := range blocks {
		_, err := a.Device.Write(ctx, block)
		if err != nil {
			return nil, fmt.Errorf("unable to write a block to HID device via interrupt OUT: %w", err)
		}
	}

	// #2: Receive ADPU response from device, in blocks
	var res Response
	var expectedSequence uint16
	for res.Length == 0 || len(res.Data) < int(res.Length) {
		data := make([]byte, a.packetSize)
		n, err := a.Device.Read(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("unable to read a block from HID device via interrupt IN: res: %v, err: %w", res, err)
		}
		if n != int(a.packetSize) {
			return nil, fmt.Errorf("read data is not equal to expected packet size, expected %d, got %d", a.packetSize, n)
		}
		res, err = a.reduceDataFrames(res, expectedSequence, data)
		if err != nil {
			return nil, fmt.Errorf("unable to reduce frame blocks, res: %v, err: %w", res, err)
		}
		expectedSequence++
	}

	a.logger.Debug("ADPU Response", "res", log.HexDisplay(res.Data))

	return res.Data, nil
}

func (a *protocolImpl) Send(ctx context.Context, cla, ins, p1, p2 uint8, data []byte) ([]byte, uint16, error) {
	var command []byte
	if len(data) > 255 {
		return nil, 0, fmt.Errorf("maximum data length of ADPU command exceeded, expected <256, got %d", len(data))
	}
	a.logger.Debug("Sending ADPU command", "cla", cla, "ins", ins, "p1", p1, "p2", p2, "data", log.HexDisplay(data))
	command = append([]byte{cla, ins, p1, p2, uint8(len(data))}, data...)

	res, err := a.Exchange(ctx, command)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to exchange ADPU, command: %s, err: %w", log.HexDisplay(command), err)
	}

	sw := binary.BigEndian.Uint16(res[len(res)-2:])

	return res[:len(res)-2], sw, nil
}

func Send[RS Unmarshaler, RQ Marshaler](ctx context.Context, proto Protocol, cla, ins, p1, p2 uint8, data RQ, res RS) error {
	dataBytes, err := Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %w", err)
	}
	response, sw, err := proto.Send(ctx, cla, ins, p1, p2, dataBytes)
	if err != nil {
		return fmt.Errorf("unable to send command to device: %w", err)
	}

	if sw != SW_OK {
		return fmt.Errorf("sw code: %s: %w", SWMessage[sw], ErrSWNotOK)
	}

	if err := Unmarshal(response, res); err != nil {
		return fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return nil
}
