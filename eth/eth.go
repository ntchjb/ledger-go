package eth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/log"
)

type EthereumApp interface {
	// Get Ledger Ethereum app's configurations
	GetConfiguration(ctx context.Context) (schema.GetConfigurationResponse, error)
	// Get address based on BIP-32 path string i.e. "m'/44'/60'/2'/0/0"
	GetAddress(ctx context.Context, bip32Path string, needHWConfirm bool, chaincode bool, chainID uint64) (schema.GetAddressResponse, error)
	// Sign a raw transaction and get signature. `rawTx` is RLP-encoded Ethereum transaction payload (EIP-155 or EIP-2718 TransactionPayload)
	SignTransaction(ctx context.Context, bip32Path string, rawTx []byte) (schema.SignTransactionResponse, error)
}

type ethereumAppImpl struct {
	proto  adpu.Protocol
	logger *slog.Logger
}

func NewEthereumApp(proto adpu.Protocol, logger *slog.Logger) EthereumApp {
	return &ethereumAppImpl{
		proto:  proto,
		logger: logger,
	}
}

func (e *ethereumAppImpl) GetConfiguration(ctx context.Context) (schema.GetConfigurationResponse, error) {
	var conf schema.GetConfigurationResponse
	if err := adpu.Send(
		ctx,
		e.proto,
		ADPU_CLA, ADPU_INS_GET_CONFIGURATION, 0x00, 0x00,
		&adpu.EmptyData{},
		&conf,
	); err != nil {
		return conf, fmt.Errorf("unable to send get configuration to device: %w", err)
	}

	return conf, nil
}

func (e *ethereumAppImpl) GetAddress(ctx context.Context, bip32Path string, needHWConfirm bool, chaincode bool, chainID uint64) (schema.GetAddressResponse, error) {
	req := schema.GetAddressRequest{
		BIP32Path: schema.BIP32Path(bip32Path),
		ChainID:   chainID,
	}
	e.logger.Debug("Get address request", "bip32Path", req.BIP32Path, "chainID", chainID)
	var res schema.GetAddressResponse

	var p1 uint8
	var p2 uint8
	if needHWConfirm {
		p1 = 0x01
	}
	if chaincode {
		p2 = 0x01
	}

	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_GET_PUBLIC_KEY, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send get address to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) SignTransaction(ctx context.Context, bip32Path string, rawTx []byte) (schema.SignTransactionResponse, error) {
	req := schema.SignTransactionRequest{
		BIP32Path: schema.BIP32Path(bip32Path),
		RawTX:     rawTx,
	}
	var res schema.SignTransactionResponse
	var resBuf []byte
	var sw uint16
	var err error

	e.logger.Debug("Sign tx request", "bip32Path", req.BIP32Path, "rawTx", log.HexDisplay(req.RawTX))

	txInfo, err := schema.DecodeTxInfo(rawTx)
	if err != nil {
		return res, fmt.Errorf("unable to decode raw tx info: %w", err)
	}
	e.logger.Debug("Tx info", "type", txInfo.TxType, "chainID", txInfo.ChainID, "chainOffset", txInfo.ChainIDOffset, "to", log.HexDisplay(txInfo.To[:]), "data", log.HexDisplay(txInfo.Data))
	if int(txInfo.TxType) >= len(schema.SupportedTxTypes) || !schema.SupportedTxTypes[txInfo.TxType] {
		return res, fmt.Errorf("unsupported transaction type: 0x%X", txInfo.TxType)
	}

	reqBuf, err := adpu.Marshal(&req)
	if err != nil {
		return res, fmt.Errorf("unable to marshal sign tx request: %w", err)
	}
	pathLength := req.BIP32Path.Len()

	for offset := 0; offset < len(reqBuf); {
		chunkSize := 255

		if offset+chunkSize > len(reqBuf) {
			chunkSize = len(reqBuf) - offset
		}

		// Workaround for a bug that ADPU chunk cannot end right after `data` field
		// for EIP-155 transaction. If it's occurred, then the Ledger device will
		// be triggered to start signing the data and ignore (chainID, 0, 0) part
		// https://github.com/LedgerHQ/app-ethereum/issues/409
		// At file /src/ethUstream.c:processTxInternal():490
		// So, we try decrease chunk size if we found that this current chunk
		// has right bound index equals to beginning of chainID data
		// to prevent from ending right at it
		if offset+chunkSize == pathLength+txInfo.ChainIDOffset && offset+chunkSize != len(reqBuf) {
			e.logger.Debug("Found chunk that cuts off chain ID, try reducing it")
			chunkSize--
		}

		e.logger.Debug("Building a chunk", "offset", offset, "chunkSize", chunkSize, "chunk", log.HexDisplay(reqBuf[offset:offset+chunkSize]))

		var p1 uint8 // P1_FIRST
		var p2 uint8 // unused
		if offset > 0 {
			p1 = 0x80 // P1_MORE
		}

		resBuf, sw, err = e.proto.Send(ctx, ADPU_CLA, ADPU_INS_SIGN_TRANSACTION, p1, p2, reqBuf[offset:offset+chunkSize])
		if err != nil {
			return res, fmt.Errorf("unable to send ADPU command to sign transaction: %w", err)
		}
		if sw != adpu.SW_OK {
			return res, fmt.Errorf("SW status: %s, %w", adpu.SWMessage[sw], adpu.ErrSWNotOK)
		}

		offset += chunkSize
	}

	// Use the last assignment of `resBuf`
	if err := adpu.Unmarshal(resBuf, &res); err != nil {
		return res, fmt.Errorf("unable to unmarshal ADPU response: %w", err)
	}

	if txInfo.TxType == schema.TX_TYPE_LEGACY {
		res.V = res.V.RecoverLegacy(txInfo.ChainID)
	}

	return res, nil
}
