package eth

import (
	"context"
	"fmt"

	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/ntchjb/ledger-go/log"
)

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
