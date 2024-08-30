package eth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
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

	// Sign typed message following EIP-712 standard
	// Signature V value can be either `27` (even), or `28` (odd)
	SignEIP712Message(ctx context.Context, bip32Path string, message eip712.Message) (schema.SignDataResponse, error)

	// Set `encodeType` data to Ledger device
	// Struct name need to be sent first, followed by struct fields
	// i.e. Mail(address from, address to, string contents) can be sent to Ledger by using following steps
	// 1. EIP712SendStructDefinition(ctx, STRUCT_COMPONENT_NAME, []byte("Mail"))
	// 2. EIP712SendStructDefinition(ctx, STRUCT_COMPONENT_FIELD, []byte{0x01, 0x02, ...})
	// For STRUCT_COMPONENT_NAME, `value` is byte array representing a UTF8 string
	// For STRUCT_COMPONENT_FIELD, `value` is marshaled data of `schema.FieldDefinition`
	EIP712SendStructDefinition(ctx context.Context, component eip712.Component, value []byte) error

	// Set `encodeData` data to Ledger device
	// For DATA_COMPONENT_ROOT, `value` is a string of data type name i.e.
	// - "EIP712Domain" for domain separator data
	// - "Mail" for primary data
	// For DATA_COMPONENT_ARRAY, `value` is length of array, stored in 1 byte
	// For DATA_COMPONENT_FIELD, `value` is serialized data that is atomic data type i.e. uint, int, bytes, etc.
	EIP712SendStructData(ctx context.Context, component eip712.Component, value []byte) error

	// Provide clear signing data to Ledger device
	// This function should be called before calling `EIP712SendStructData`
	// It is usually called after `EIP712SendStructDefinition` was called
	EIP712SendClearSigningData(ctx context.Context, action eip712.Action, value []byte) error

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
	// Provide ERC20 information to be displayed during transaction signing
	// This function shall be run before `SignTransaction`
	// `info` is ERC20 information, which can be obtained from Ledger Live API
	ProvideERC20Information(ctx context.Context, info []byte) (schema.ProvideERC20InfoResponse, error)
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

func (e *ethereumAppImpl) EIP712SendStructDefinition(ctx context.Context, component eip712.Component, value []byte) error {
	req := schema.RawRequest(value)
	var res schema.EmptyResponse
	p1 := uint8(0x00)
	p2 := uint8(component)

	e.logger.Debug("Send EIP712 struct definition", "component", component, "value", log.HexDisplay(value))
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_EIP712_SEND_STRUCT_DEF, p1, p2, &req, &res); err != nil {
		return fmt.Errorf("unable to send a send struct definition command to device: %w", err)
	}

	return nil
}

func (e *ethereumAppImpl) EIP712SendClearSigningData(ctx context.Context, action eip712.Action, value []byte) error {
	req := schema.RawRequest(value)
	var res schema.EmptyResponse
	p1 := uint8(0x00)
	p2 := uint8(action)

	e.logger.Debug("Provide EIP712 clear signing data", "action", action, "value", log.HexDisplay(value))
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_EIP712_CLEAR_SIGNING, p1, p2, &req, &res); err != nil {
		return fmt.Errorf("unable to send EIP712 clear signing command to device: %w", err)
	}

	return nil
}

func (e *ethereumAppImpl) EIP712SendStructData(ctx context.Context, component eip712.Component, value []byte) error {
	var req schema.RawRequest
	var res schema.EmptyResponse

	for offset := 0; offset < len(value); {
		chunkSize := 255
		p1 := P1_PARTIAL
		p2 := uint8(component)
		if offset+chunkSize >= len(value) {
			p1 = P1_COMPLETE
			chunkSize = len(value) - offset
		}
		req = value[offset : offset+chunkSize]

		e.logger.Debug("Send EIP712 data", "val", log.HexDisplay(req), "component", component, "p1", p1, "p2", p2)
		if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_EIP712_SEND_STRUCT_DATA, p1, p2, &req, &res); err != nil {
			return fmt.Errorf("unable to send a send EIP712 data command to device: %w", err)
		}

		offset += chunkSize
	}

	return nil
}

func (e *ethereumAppImpl) sendEIP712Data(ctx context.Context, cs eip712.ClearSigning, domain eip712.Domain, coinRefRegistered map[int]uint8) eip712.WalkReader {
	return func(path string, item eip712.Item) error {
		// #1: Provide Clear signing information to device, if enabled
		if cs.Enabled && item.Type() == eip712.DATA_COMPONENT_ATOMIC {
			fieldInfo, fieldExists := cs.Fields[path]
			if !fieldExists {
				goto endOfClearSigning
			}
			e.logger.Debug("Setup ERC20 clear signing data")

			// #1.1: Provide ERC20 info based on coinRef
			if _, isERC20TokenProvided := coinRefRegistered[fieldInfo.CoinRef]; fieldInfo.Format == eip712.CSIGN_FIELD_FORMAT_TOKEN && fieldInfo.CoinRef >= 0 && !isERC20TokenProvided {
				address, ok := cs.CoinRefMap[fieldInfo.CoinRef]
				if !ok {
					return fmt.Errorf("unable to find token by coin ref: %d, coinRef: %+v", fieldInfo.CoinRef, cs.CoinRefMap)
				}
				if tokenInfo, ok := cs.ERC20Signatures.FindByChainIDAndAddress(domain.ChainID, address); ok {
					res, err := e.ProvideERC20Information(ctx, tokenInfo.Raw)
					if err != nil {
						return fmt.Errorf("unable to provide ERC20 info, contractAddress: 0x%x, err: %w", tokenInfo.ContractAddress, err)
					}

					coinRefRegistered[fieldInfo.CoinRef] = uint8(res)
				}
			}

			// #1.2: Provide ERC20 info of verifying contract address, if any (coinRef = 255 means it's verifying contract)
			if fieldInfo.Format == eip712.CSIGN_FIELD_FORMAT_AMOUNT && fieldInfo.CoinRef == 255 {
				address := cs.CoinRefMap[255]

				if tokenInfo, ok := cs.ERC20Signatures.FindByChainIDAndAddress(domain.ChainID, address); ok {
					if _, err := e.ProvideERC20Information(ctx, tokenInfo.Raw); err != nil {
						return fmt.Errorf("unable to provide ERC20 info, contractAddress: 0x%x, err: %w", tokenInfo.ContractAddress, err)
					}

					coinRefRegistered[fieldInfo.CoinRef] = 255
				}
			}

			// #1.3: Provide EIP712 clear signing data i.e. display name of the atomic field, based on field info
			eip712CSignPayload, err := fieldInfo.Payload(coinRefRegistered)
			if err != nil {
				return fmt.Errorf("unable to create EIP712 payload for clear signing field: %w", err)
			}
			eip712CSignAction, err := fieldInfo.Action()
			if err != nil {
				return fmt.Errorf("cannot get action from EIP712 field, format: %s, err: %w", fieldInfo.Format, err)
			}
			if err := e.EIP712SendClearSigningData(ctx, eip712CSignAction, eip712CSignPayload); err != nil {
				return fmt.Errorf("unable to send EIP712 clear signing data: %w", err)
			}
		}

	endOfClearSigning:
		// #2: Send EIP712 data to device
		if item.Type() == eip712.DATA_COMPONENT_ROOT {
			return nil
		}

		dataCmd := item.DataCommand()
		if err := e.EIP712SendStructData(ctx, dataCmd.Component, dataCmd.Value); err != nil {
			return fmt.Errorf("unable to send data: %w", err)
		}

		return nil
	}
}

func (e *ethereumAppImpl) SignEIP712Message(ctx context.Context, bip32Path string, message eip712.Message) (schema.SignDataResponse, error) {
	var res schema.SignDataResponse
	// #1: Send type definition
	for _, typeDef := range message.Types {
		if err := e.EIP712SendStructDefinition(ctx, eip712.TYPE_COMPONENT_NAME, []byte(typeDef.Name)); err != nil {
			return res, fmt.Errorf("unable to send EIP712 struct definition, type name: %w", err)
		}
		for _, member := range typeDef.Members {
			structDefBytes, err := member.MarshalADPU()
			if err != nil {
				return res, fmt.Errorf("unable to marshal FieldDefinition: %w", err)
			}
			if err := e.EIP712SendStructDefinition(ctx, eip712.TYPE_COMPONENT_FIELD, structDefBytes); err != nil {
				return res, fmt.Errorf("unable to send EIP712 struct definition, field type: %w", err)
			}
		}
	}

	// #2: Activate clear signing, if enabled
	if message.ClearSigning.Enabled {
		e.logger.Debug("Activate clear signing for EIP712")
		if err := e.EIP712SendClearSigningData(ctx, eip712.ACTION_ACTIVATE, nil); err != nil {
			return res, fmt.Errorf("unable to activate clear signing: %w", err)
		}
	}

	// #3: Send domain root type and data
	domainStructItem := message.Domain.StructItem()
	coinRefRegisteredOnDevice := make(map[int]uint8)
	domainRootCmd := domainStructItem.DataCommand()
	if err := e.EIP712SendStructData(ctx, domainRootCmd.Component, domainRootCmd.Value); err != nil {
		return res, fmt.Errorf("unable to set domain root data: %w", err)
	}
	if err := domainStructItem.Walk("", e.sendEIP712Data(ctx, eip712.ClearSigning{}, message.Domain, coinRefRegisteredOnDevice)); err != nil {
		return res, fmt.Errorf("unable to send domain data: %w", err)
	}

	// #4: Send contract name as clear signing data, if any
	if message.ClearSigning.Enabled {
		if err := message.SetCoinRefMap(message.Primary); err != nil {
			return res, fmt.Errorf("unable to set coin ref map for primary data: %w", err)
		}
		payload := message.ClearSigning.ContractPayload()
		if err := e.EIP712SendClearSigningData(ctx, eip712.ACTION_MESSAGE_INFO, payload); err != nil {
			return res, fmt.Errorf("unable to send EIP712 clear signing data; contract info: %w", err)
		}
	}

	// #5: Send primary data
	primaryRootCmd := message.Primary.DataCommand()
	if err := e.EIP712SendStructData(ctx, primaryRootCmd.Component, primaryRootCmd.Value); err != nil {
		return res, fmt.Errorf("unable to set primary root data: %w", err)
	}
	if err := message.Primary.Walk("", e.sendEIP712Data(ctx, message.ClearSigning, message.Domain, coinRefRegisteredOnDevice)); err != nil {
		return res, fmt.Errorf("unable to send primary data: %w", err)
	}

	// #6: Send HD wallet path as the last command and return signature
	p1, p2 := uint8(0x00), uint8(0x01)
	req := schema.BIP32Path(bip32Path)
	e.logger.Debug("Sign EIP712 message", "bip32Path", bip32Path)
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_SIGN_EIP712, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send sign EIP712 command to device: %w", err)
	}

	return res, nil
}

func (e *ethereumAppImpl) ProvideERC20Information(ctx context.Context, info []byte) (schema.ProvideERC20InfoResponse, error) {
	req := schema.RawRequest(info)
	var res schema.ProvideERC20InfoResponse
	var p1, p2 uint8

	e.logger.Debug("Provide ERC20 information", "info", log.HexDisplay(info))
	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_PROVIDE_ERC20_INFO, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send provide ERC20 information command to device: %w", err)
	}

	return res, nil
}
