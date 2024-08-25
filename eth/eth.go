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
	SignTransaction(ctx context.Context, bip32Path string, rawTx []byte) (schema.SignDataResponse, error)
	// Sign a personal message following ERC-191 standard
	// The message is usually a string, but it supports arbitrary data
	// Signature V value can be either `27` (even), or `28` (odd)
	SignPersonalMessage(ctx context.Context, bip32Path string, message []byte) (schema.SignDataResponse, error)

	// // Sign typed message following EIP-712 standard
	// // Signature V value can be either `27` (even), or `28` (odd)
	// SignEIP712Message(ctx context.Context, bip32Path string, domainSeparator eip712.Message) (schema.SignDataResponse, error)

	// Sign typed message (hashed format) following EIP-712 standard
	// Signature V value can be either `27` (even), or `28` (odd)
	SignEIP712MessageHash(ctx context.Context, bip32Path string, domainSeparatorHash []byte, messageHash []byte) (schema.SignDataResponse, error)

	// Get BLS12-381 public key following EIP-2333 standard by given BIP-32 path
	// BIP-32 path follows EIP-2334 standard i.e.
	// - m/12381/3600/0/0 (for withdrawal key)
	// - m/12381/3600/0/0/0 (for signing key)
	ETH2GetPublicKey(ctx context.Context, bip32Path string, needHWConfirm bool) (schema.ETH2PublicKey, error)
	// Set index of withdrawal key used as withdrawal credential with ETH2 deposit contract call
	// BIP-32 path follows EIP-2334 standard i.e. m/12381/3600/0/0
	// This function shall be run before `SignTransaction` when signing ETH2 deposit transaction
	// Otherwise, withdrawal index 0 is used
	ETH2SetWithdrawalIndex(ctx context.Context, index uint32) error

	// Get public key of Curve25519 key pair to perform end-to-end encryption (X25519)
	GetPrivacyPublicKey(ctx context.Context, bip32Path string, needHWConfirm bool) (schema.GetPrivacyPublicKeyResponse, error)

	// Get shared key of X25519
	GetPrivacySharedSecret(ctx context.Context, bip32Path string, remotePublicKey []byte, needHWConfirm bool) (schema.GetPrivacySharedSecretResponse, error)

	// // Get 4-byte challenge data from Ledger device to be signed by trusted Ledger Live Server, used by clear signing
	// // Currently, it's used for displaying domain name instead of `to` address during transaction signing.
	// // This function shall be run before `ProvideDomainName`
	// GetChallenge(ctx context.Context) (schema.Challenge, error)
	// // Provide domain name i.e. ENS to be displayed during transaction signing in place of `to` address.
	// // This function shall be run before `SignTransaction`
	// // `info` is TVL data (tag-value-length) that can be obtained from Ledger Live API
	// ProvideDomainNameInformation(ctx context.Context, info []byte) error
	// // Provide NFT information to be displayed during transaction signing
	// // This function shall be run before `SignTransaction`
	// // `info` is NFT information, which can be obtained from Ledger Live API
	// ProvideNFTInformation(ctx context.Context, info []byte) error
	// // Provide ERC20 information to be displayed during transaction signing
	// // This function shall be run before `SignTransaction`
	// // `info` is ERC20 information, which can be obtained from Ledger Live API
	// ProvideERC20Information(ctx context.Context, info []byte) error
	// // Provide name of a plugin to interpret contract data, used by clear signing.
	// // The plugin determines contract address and its method selectors (contract function that is called)
	// // and provide information on Ledger device display during transaction signing
	// // This function shall be run before `SignTransaction`
	// // `info` can be obtained from Ledger Live API
	// SetPlugin(ctx context.Context, info []byte) error
	// // Provide name of an external plugin to interpret contract data, used by clear signing.
	// // The plugin determines contract address and its method selectors (contract function that is called)
	// // and provide information on Ledger device display during transaction signing
	// // This function shall be run before `SignTransaction`
	// // `info` can be obtained from Ledger Live API
	// SetExternalPlugin(ctx context.Context, info []byte) error
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
	var res schema.GetAddressResponse

	p1, p2 := P1_WITHOUT_CONFIRM, uint8(0x00)
	if needHWConfirm {
		p1 = P1_WITH_CONFIRM
	}
	if chaincode {
		p2 = 0x01
	}

	e.logger.Debug("Get address request", "bip32Path", req.BIP32Path, "chainID", chainID)

	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_GET_PUBLIC_KEY, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send get address to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) SignTransaction(ctx context.Context, bip32Path string, rawTx []byte) (schema.SignDataResponse, error) {
	req := schema.SignTxRequest{
		BIP32Path: schema.BIP32Path(bip32Path),
		Data:      rawTx,
	}
	var res schema.SignDataResponse
	var resBuf []byte
	var sw uint16
	var err error

	e.logger.Debug("Sign tx request", "bip32Path", req.BIP32Path, "rawTx", log.HexDisplay(req.Data))

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

		p1 := P1_FIRST_CHUNK
		var p2 uint8 // unused
		if offset > 0 {
			p1 = P1_MORE_CHUNK
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

func (e *ethereumAppImpl) SignPersonalMessage(ctx context.Context, bip32Path string, message []byte) (schema.SignDataResponse, error) {
	req := schema.SignPersonalMessageRequest{
		BIP32Path: schema.BIP32Path(bip32Path),
		Data:      message,
	}
	var res schema.SignDataResponse
	var err error
	var resBuf []byte
	var sw uint16

	e.logger.Debug("Sign personal message", "bip32Path", req.BIP32Path, "message", log.HexDisplay(message))

	reqBuf, err := adpu.Marshal(&req)
	if err != nil {
		return res, fmt.Errorf("unable to marshal sign personal message request: %w", err)
	}

	for offset := 0; offset < len(reqBuf); {
		chunkSize := 255
		if offset+chunkSize > len(reqBuf) {
			chunkSize = len(reqBuf) - offset
		}
		p1 := P1_FIRST_CHUNK
		var p2 uint8 // unused
		if offset > 0 {
			p1 = P1_MORE_CHUNK
		}
		e.logger.Debug("Building a chunk", "offset", offset, "chunkSize", chunkSize, "chunk", log.HexDisplay(reqBuf[offset:offset+chunkSize]))

		resBuf, sw, err = e.proto.Send(ctx, ADPU_CLA, ADPU_INS_SIGN_PERSONAL_MESSAGE, p1, p2, reqBuf[offset:offset+chunkSize])
		if err != nil {
			return res, fmt.Errorf("unable to send ADPU command to sign personal message: %w", err)
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

	return res, nil
}

func (e *ethereumAppImpl) SignEIP712MessageHash(ctx context.Context, bip32Path string, domainSeparatorHash []byte, messageHash []byte) (schema.SignDataResponse, error) {
	req := schema.SignEIP712HashedRequest{
		BIP32Path:             schema.BIP32Path(bip32Path),
		HashedDomainSeparator: [32]byte(domainSeparatorHash),
		HashedMessage:         [32]byte(messageHash),
	}
	var res schema.SignDataResponse
	var p1, p2 uint8

	if len(domainSeparatorHash) != 32 {
		return res, fmt.Errorf("unexpeted hashed domain separator length, expected 32, got %d", len(domainSeparatorHash))
	}
	if len(messageHash) != 32 {
		return res, fmt.Errorf("unexpected hashed domain separator length, expected 32, got %d", len(messageHash))
	}

	e.logger.Debug("Sign EIP712 message", "bip32Path", req.BIP32Path, "domain", log.HexDisplay(req.HashedDomainSeparator[:]), "message", log.HexDisplay(req.HashedMessage[:]))
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_SIGN_EIP712, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send sign EIP712 command to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) ETH2GetPublicKey(ctx context.Context, bip32Path string, needHWConfirm bool) (schema.ETH2PublicKey, error) {
	req := schema.BIP32Path(bip32Path)
	var res schema.ETH2PublicKey
	p1, p2 := P1_WITHOUT_CONFIRM, uint8(0x00)
	if needHWConfirm {
		p1 = P1_WITH_CONFIRM
	}

	e.logger.Debug("Get ETH2 public key", "bip32Path", bip32Path, "confirm", needHWConfirm)
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_ETH2_GET_PUBLIC_KEY, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send ETH2 get public key command to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) ETH2SetWithdrawalIndex(ctx context.Context, index uint32) error {
	req := schema.ETH2WithdrawalIndex(index)
	var res schema.EmptyResponse
	var p1, p2 uint8

	e.logger.Debug("Set ETH2 withdrawal index", "index", index)
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_ETH2_SET_WITHDRAWAL_INDEX, p1, p2, &req, &res); err != nil {
		return fmt.Errorf("unable to send ETH2 set withdrawal index command to device: %w", err)
	}

	return nil
}

func (e *ethereumAppImpl) GetPrivacyPublicKey(ctx context.Context, bip32Path string, needHWConfirm bool) (schema.GetPrivacyPublicKeyResponse, error) {
	req := schema.BIP32Path(bip32Path)
	var res schema.GetPrivacyPublicKeyResponse
	p1, p2 := P1_WITHOUT_CONFIRM, uint8(0x00)
	if needHWConfirm {
		p1 = P1_WITH_CONFIRM
	}

	e.logger.Debug("Get privacy public key", "bip32Path", bip32Path, "confirm", needHWConfirm)
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_PRIVACY_OPERATION, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send get privacy public key command to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) GetPrivacySharedSecret(ctx context.Context, bip32Path string, remotePublicKey []byte, needHWConfirm bool) (schema.GetPrivacySharedSecretResponse, error) {
	req := schema.GetPrivacySharedSecretRequest{
		Path:            schema.BIP32Path(bip32Path),
		RemotePublicKey: remotePublicKey,
	}
	var res schema.GetPrivacySharedSecretResponse
	p1, p2 := P1_WITHOUT_CONFIRM, uint8(0x01)
	if needHWConfirm {
		p1 = P1_WITH_CONFIRM
	}

	e.logger.Debug("Get shared secret key", "bip32Path", bip32Path, "confirm", needHWConfirm, "remotePublicKey", log.HexDisplay(remotePublicKey))
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_PRIVACY_OPERATION, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send get shared secret command to device: %w", err)
	}

	return res, nil
}
