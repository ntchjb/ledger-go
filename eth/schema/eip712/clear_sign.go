package eip712

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/ntchjb/ledger-go/eth/schema"
)

var (
	ErrRegisteredCoinRefNotFound = errors.New("registered coin ref not found")
	ErrUnknownFieldFormat        = errors.New("unknown clear signing field format")
)

type CSignContract struct {
	Label string
	// Signature returned from Ledger Live server is hexadecimal string.
	// It needs to be converted to byte array
	Signature []byte
}

func (c *CSignContract) DisplayPayload() []byte {
	length := byte(len(c.Label))

	return append([]byte{length}, []byte(c.Label)...)
}

func (c *CSignContract) SignaturePayload() []byte {
	length := byte(len(c.Signature))

	return append([]byte{length}, c.Signature...)
}

type CSignFieldFormat string

const (
	CSIGN_FIELD_FORMAT_RAW      CSignFieldFormat = "raw"
	CSIGN_FIELD_FORMAT_TOKEN    CSignFieldFormat = "token"
	CSIGN_FIELD_FORMAT_AMOUNT   CSignFieldFormat = "amount"
	CSIGN_FIELD_FORMAT_DATETIME CSignFieldFormat = "datetime"
)

type CSignField struct {
	Format    CSignFieldFormat
	Label     string
	Signature []byte
	CoinRef   int
}

func (c *CSignField) DisplayPayload() []byte {
	length := byte(len(c.Label))

	return append([]byte{length}, []byte(c.Label)...)
}

func (c *CSignField) SignaturePayload() []byte {
	length := byte(len(c.Signature))

	return append([]byte{length}, c.Signature...)
}

func (c *CSignField) Payload(registeredCoinRef map[int]uint8) ([]byte, error) {
	display, signature := c.DisplayPayload(), c.SignaturePayload()
	switch c.Format {
	case CSIGN_FIELD_FORMAT_RAW:
		fallthrough
	case CSIGN_FIELD_FORMAT_DATETIME:
		return append(display, signature...), nil
	case CSIGN_FIELD_FORMAT_TOKEN:
		if idx, ok := registeredCoinRef[c.CoinRef]; !ok {
			return nil, ErrRegisteredCoinRefNotFound
		} else {
			return append([]byte{idx}, signature...), nil
		}
	case CSIGN_FIELD_FORMAT_AMOUNT:
		if idx, ok := registeredCoinRef[c.CoinRef]; !ok {
			return nil, ErrRegisteredCoinRefNotFound
		} else {
			return append(append(display, idx), signature...), nil
		}
	default:
		return nil, ErrUnknownFieldFormat
	}
}

func (c *CSignField) Action() (Action, error) {
	switch c.Format {
	case CSIGN_FIELD_FORMAT_RAW:
		return ACTION_RAW, nil
	case CSIGN_FIELD_FORMAT_DATETIME:
		return ACTION_DATETIME, nil
	case CSIGN_FIELD_FORMAT_TOKEN:
		return ACTION_AMOUNT_TOKEN_JOIN, nil
	case CSIGN_FIELD_FORMAT_AMOUNT:
		return ACTION_AMOUNT_VALUE_JOIN, nil
	default:
		return 0, ErrUnknownFieldFormat
	}
}

type CSignTokenInfo struct {
	ContractAddress schema.Address
	Ticker          string
	Decimals        uint32
	ChainID         uint32
	Signature       []byte

	Raw []byte
}

type ERC20Signatures map[[24]byte]CSignTokenInfo

func (e ERC20Signatures) FindByChainIDAndAddress(chainID *uint256.Int, address schema.Address) (CSignTokenInfo, bool) {
	var key [24]byte
	bChainID := chainID.Bytes32()
	copy(key[:4], bChainID[len(bChainID)-4:])
	copy(key[4:], address[:])

	res, ok := e[key]
	return res, ok
}

func ParseERC20SignatureBlobs(blob string) (ERC20Signatures, error) {
	res := make(map[[24]byte]CSignTokenInfo)
	blobBytes, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return res, fmt.Errorf("unable to decode base64 blob: %w", err)
	}
	for i := 0; i < len(blobBytes); {
		var tokenInfo CSignTokenInfo
		var key [24]byte

		recordLength := binary.BigEndian.Uint32(blobBytes[i : i+4])
		i += 4

		record := blobBytes[i : i+int(recordLength)]

		recIdx := 0
		tickerLength := int(record[recIdx])
		recIdx += 1

		tokenInfo.Ticker = string(record[recIdx : recIdx+tickerLength])
		recIdx += tickerLength

		copy(tokenInfo.ContractAddress[:], record[recIdx:recIdx+20])
		copy(key[4:], record[recIdx:recIdx+20])
		recIdx += 20

		tokenInfo.Decimals = binary.BigEndian.Uint32(record[recIdx : recIdx+4])
		recIdx += 4

		tokenInfo.ChainID = binary.BigEndian.Uint32(record[recIdx : recIdx+4])
		copy(key[:4], record[recIdx:recIdx+4])
		recIdx += 4

		tokenInfo.Signature = record[recIdx:]
		tokenInfo.Raw = record

		res[key] = tokenInfo

		i += int(recordLength)
	}

	return res, nil
}

type ClearSigning struct {
	Enabled bool

	ContractInfo CSignContract
	// Fields' key is path
	Fields          map[string]CSignField
	ERC20Signatures ERC20Signatures
	// Key is CoinRef and Value is path
	// If Key is 255, it is domain's verifying contract
	CoinRefMap map[int]schema.Address
}

func (c *ClearSigning) ContractPayload() []byte {
	fieldCount := byte(len(c.Fields))
	return append(append(c.ContractInfo.DisplayPayload(), fieldCount), c.ContractInfo.SignaturePayload()...)
}
