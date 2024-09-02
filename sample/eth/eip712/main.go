package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/holiman/uint256"
	"github.com/ntchjb/gohid/hid"
	"github.com/ntchjb/gohid/manager"
	"github.com/ntchjb/gohid/usb"
	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/device"
	"github.com/ntchjb/ledger-go/eth"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/ntchjb/ledger-go/log"
	"github.com/ntchjb/ledger-go/sample/eth/util"
)

func main() {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	usbCtx := usb.NewGOUSBContext()
	man := manager.NewDeviceManager(usbCtx, logger)

	defer func() {
		if err := man.Close(); err != nil {
			logger.Error("unable to close device manager", "err", err)
		}
	}()

	// Open Ledger Nano S
	// hidDevice, err := man.Open(0x2C97, 0x1011, hid.DeviceConfig{
	// 	StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	// })

	// Open Ledger Nano S Plus
	hidDevice, err := man.Open(0x2C97, 0x5011, hid.DeviceConfig{
		StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	})

	if err != nil {
		logger.Error("unable to open device", "err", err)
		return
	}
	defer hidDevice.Close()

	if err := hidDevice.SetAutoDetach(true); err != nil {
		logger.Error("unable to set auto detach", "err", err)
		return
	}
	if err := hidDevice.SetTarget(1, 0, 0); err != nil {
		logger.Error("unable to set target of hid device", "err", err)
		return
	}

	// Get report descriptor
	desc, err := hidDevice.GetReportDescriptor()
	if err != nil {
		logger.Error("unable to get report descriptor", "err", err)
		return
	}
	str, err := desc.String()
	if err != nil {
		logger.Error("unable to get descriptor string", "err", err)
		return
	}
	fmt.Printf("%s\n", str)
	fmt.Printf("%v\n", desc)

	// Create new Ethereum app instance
	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)
	ethApp := eth.NewEthereumApp(adpuProto, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	walletPath := "m'/44'/60'/0'/0/0"

	// DATA from https://crypto-assets-service.api.ledger.com/v1/dapps?output=eip712_signatures&eip712_signatures_version=v2&chain_id=10&contracts=0x000000000022d473030f116ddee9f6b43ac78ba3
	uniswapXLimitOrderContractSignature, _ := hex.DecodeString("3045022100f607f91959ba77569e1bbc520fd61ebd0cf2c6b0b4bfa449c45e86ac49f048e602200a1f105838d380ef60f765dcb0d3bcfd2eb9af8dee82994a942bf804eb5c144c")
	uniXApproveToSpenderSig, _ := hex.DecodeString("304402203ae7648a1fcc87edd672587dcd9c4222aef9b119eb5573945982eb4763c9c110022072d0a4d1e23db36c3b4852bc61b8500e0a9b4a58d56ed6b71d8491e154e1773d")
	uniXApproveTokenSig, _ := hex.DecodeString("3045022100d89ed36285b1474f6caac45467ccf5ded7e63218542cb36cbbc25970416479370220296bb6d4643dd43d842c0f52227fc3497c23f8402404a50537e8e6e76a0406a0")
	uniXApproveAmountSig, _ := hex.DecodeString("304402201e0da0f02cca490ca1c231089ef95664fa830ffa1225e1d66aa217034f988d7b02202fb83a698424fb3434ec61cfeb6db7ac565ea318145450544b6a3d509682f96b")
	uniXToSwapTokenSig, _ := hex.DecodeString("304502210090e29b4ae8364ce6fdf0a1162a381baf1db0d9654e4098e98aea191bf5dda392022014e87bb5261fb8ab9d1d1694ed928fbadfa81810fafffe5b684d255c4570ee1d")
	uniXToSwapAmountSig, _ := hex.DecodeString("3044022075f4050f8ccd04f0832ac81a5c73d12ddd78baad003e81f5931ce2f43303f14402203ac51a3456ce84ad7c934fe30a469b6874d47510e4b097b386aff5faa214b975")
	uniXTokenToReceiveSig, _ := hex.DecodeString("3044022053bc0c1caba1f2a589ced91e416486419aa499e625d8fb4256675a3216bec772022057698f1ed49eb612601479aaa33ab77b635ab38dcce54f8d354e46f08a36a566")
	uniXMinAmountToReceiveSig, _ := hex.DecodeString("3045022100f1748b0339fccd0dc2e7780d701816b551b92c01c9a582387c9c5f19310c4d48022070a3ab6e0d49b285ca87f58ccb4eeccc979389382ffd6390e0d0398771cd3cff")
	uniXOnAddressSig, _ := hex.DecodeString("3045022100cd701a6cf3d4150d9ac6efd79e72f790772433dbde62cf4b537b5ae2c51e0d44022009372e93db760ff9d6fe88c9a912d1e1595fe0fa85aa53ef759e13ccf95ca87f")
	uniXApprovalExpire, _ := hex.DecodeString("3044022018740d5b88a5a9245b59148cfb26c2728af523a4ffe23329646c6f07454721c90220426efe50d47b3f6f051ff70a132d93d3d549dd2b9823725bc2fd8e8affaf1dc7")

	domain := eip712.Domain{
		Name:    "Permit2",
		ChainID: uint256.NewInt(10),
		VerifyingContract: schema.Address{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0xD4, 0x73, 0x03, 0x0F,
			0x11, 0x6d, 0xDE, 0xE9, 0xF6, 0xB4, 0x3a, 0xC7, 0x8B, 0xA3,
		},
	}
	message := eip712.Message{
		Types: eip712.TypeStructs{
			domain.TypeStruct(),

			{
				Name: "LimitOrder",
				Members: []eip712.FieldDefinition{
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_CUSTOM,
						},
						CustomTypeName: "OrderInfo",
						KeyName:        "info",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "inputToken",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "inputAmount",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:    eip712.FIELD_TYPE_DESC_TYPE_CUSTOM,
							IsArray: true,
						},
						CustomTypeName: "OutputToken",
						ArrayLevels: []eip712.FieldArrayLevel{
							{
								Type: eip712.STRUCT_DEF_ARRAY_TYPE_DYNAMIC,
							},
						},
						KeyName: "outputs",
					},
				},
			},
			{
				Name: "OrderInfo",
				Members: []eip712.FieldDefinition{
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "reactor",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "swapper",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "nonce",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "deadline",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "additionalValidationContract",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_DYNAMIC_SIZED_BYTES,
						},
						KeyName: "additionalValidationData",
					},
				},
			},
			{
				Name: "OutputToken",
				Members: []eip712.FieldDefinition{
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "token",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "amount",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "recipient",
					},
				},
			},
			{
				Name: "PermitWitnessTransferFrom",
				Members: []eip712.FieldDefinition{
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_CUSTOM,
						},
						CustomTypeName: "TokenPermissions",
						KeyName:        "permitted",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "spender",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "nonce",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "deadline",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_CUSTOM,
						},
						CustomTypeName: "LimitOrder",
						KeyName:        "witness",
					},
				},
			},
			{
				Name: "TokenPermissions",
				Members: []eip712.FieldDefinition{
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type: eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "token",
					},
					{
						TypeDescription: eip712.FieldTypeDescription{
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
							IsSizeSpecified: true,
						},
						TypeSize: 32,
						KeyName:  "amount",
					},
				},
			},
		},
		Domain: domain,
		Primary: eip712.StructItem{
			TypeName: "PermitWitnessTransferFrom",
			Members: []eip712.StructItemMember{
				{
					Name: "permitted",
					Item: eip712.StructItem{
						TypeName: "TokenPermissions",
						Members: []eip712.StructItemMember{
							{
								Name: "token",
								Item: eip712.AtomicItem{
									Item: eip712.AddressData{
										0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
										0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42,
									},
								},
							},
							{
								Name: "amount",
								Item: eip712.AtomicItem{
									Item: eip712.NumberData{
										Num:     uint256.MustFromDecimal("100000000000000000000"),
										NumBits: 256,
									},
								},
							},
						},
					},
				},
				{
					Name: "spender",
					Item: eip712.AtomicItem{
						Item: eip712.AddressData{
							0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
							0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
						},
					},
				},
				{
					Name: "nonce",
					Item: eip712.AtomicItem{
						Item: eip712.NumberData{
							Num:     uint256.NewInt(124),
							NumBits: 256,
						},
					},
				},
				{
					Name: "deadline",
					Item: eip712.AtomicItem{
						Item: eip712.NumberData{
							Num:     uint256.NewInt(1724998503),
							NumBits: 256,
						},
					},
				},
				{
					Name: "witness",
					Item: eip712.StructItem{
						TypeName: "LimitOrder",
						Members: []eip712.StructItemMember{
							{
								Name: "info",
								Item: eip712.StructItem{
									TypeName: "OrderInfo",
									Members: []eip712.StructItemMember{
										{
											Name: "reactor",
											Item: eip712.AtomicItem{
												Item: eip712.AddressData{
													0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
													0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
												},
											},
										},
										{
											Name: "swapper",
											Item: eip712.AtomicItem{
												Item: eip712.AddressData{
													0xF1, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
													0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
												},
											},
										},
										{
											Name: "nonce",
											Item: eip712.AtomicItem{
												Item: eip712.NumberData{
													Num:     uint256.NewInt(435),
													NumBits: 256,
												},
											},
										},
										{
											Name: "deadline",
											Item: eip712.AtomicItem{
												Item: eip712.NumberData{
													Num:     uint256.NewInt(1724998883),
													NumBits: 256,
												},
											},
										},
										{
											Name: "additionalValidationContract",
											Item: eip712.AtomicItem{
												Item: eip712.AddressData{
													0xE1, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
													0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
												},
											},
										},
										{
											Name: "additionalValidationData",
											Item: eip712.AtomicItem{
												Item: eip712.AddressData{
													0xE2, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
													0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
												},
											},
										},
									},
								},
							},
							{
								Name: "inputToken",
								Item: eip712.AtomicItem{
									Item: eip712.AddressData{
										0x94, 0xb0, 0x08, 0xaA, 0x00, 0x57, 0x9c, 0x13, 0x07, 0xB0,
										0xEF, 0x2c, 0x49, 0x9a, 0xD9, 0x8a, 0x8c, 0xe5, 0x8e, 0x58,
									},
								},
							},
							{
								Name: "inputAmount",
								Item: eip712.AtomicItem{
									Item: eip712.NumberData{
										Num:     uint256.MustFromDecimal("300000000"),
										NumBits: 256,
									},
								},
							},
							{
								Name: "outputs",
								Item: eip712.ArrayItem{
									eip712.StructItem{
										TypeName: "OutputToken",
										Members: []eip712.StructItemMember{
											{
												Name: "token",
												Item: eip712.AtomicItem{
													Item: eip712.AddressData{
														0x0b, 0x2C, 0x63, 0x9c, 0x53, 0x38, 0x13, 0xf4, 0xAa, 0x9D,
														0x78, 0x37, 0xCA, 0xf6, 0x26, 0x53, 0xd0, 0x97, 0xFf, 0x85,
													},
												},
											},
											{
												Name: "amount",
												Item: eip712.AtomicItem{
													Item: eip712.NumberData{
														Num:     uint256.MustFromDecimal("200000000"),
														NumBits: 256,
													},
												},
											},
											{
												Name: "recipient",
												Item: eip712.AtomicItem{
													Item: eip712.AddressData{
														0x55, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
														0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
													},
												},
											},
										},
									},
									eip712.StructItem{
										TypeName: "OutputToken",
										Members: []eip712.StructItemMember{
											{
												Name: "token",
												Item: eip712.AtomicItem{
													Item: eip712.AddressData{
														0x0b, 0x2C, 0x63, 0x9c, 0x53, 0x38, 0x13, 0xf4, 0xAa, 0x9D,
														0x78, 0x37, 0xCA, 0xf6, 0x26, 0x53, 0xd0, 0x97, 0xFf, 0x85,
													},
												},
											},
											{
												Name: "amount",
												Item: eip712.AtomicItem{
													Item: eip712.NumberData{
														Num:     uint256.MustFromDecimal("100000000"),
														NumBits: 256,
													},
												},
											},
											{
												Name: "recipient",
												Item: eip712.AtomicItem{
													Item: eip712.AddressData{
														0x56, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
														0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ClearSigning: eip712.ClearSigning{
			Enabled: true,
			ContractInfo: eip712.CSignContract{
				Label:     "UniswapX Limit Order",
				Signature: uniswapXLimitOrderContractSignature,
			},
			Fields: map[string]eip712.CSignField{
				"spender": {
					Format:    eip712.CSIGN_FIELD_FORMAT_RAW,
					Label:     "Approve to spender",
					Signature: uniXApproveToSpenderSig,
				},
				"permitted.token": {
					Format:    eip712.CSIGN_FIELD_FORMAT_TOKEN,
					Label:     "Approve amount",
					Signature: uniXApproveTokenSig,
					CoinRef:   0,
				},
				"permitted.amount": {
					Format:    eip712.CSIGN_FIELD_FORMAT_AMOUNT,
					Label:     "Approve amount",
					Signature: uniXApproveAmountSig,
					CoinRef:   0,
				},
				"witness.inputToken": {
					Format:    eip712.CSIGN_FIELD_FORMAT_TOKEN,
					Label:     "To swap",
					Signature: uniXToSwapTokenSig,
					CoinRef:   1,
				},
				"witness.inputAmount": {
					Format:    eip712.CSIGN_FIELD_FORMAT_AMOUNT,
					Label:     "To swap",
					Signature: uniXToSwapAmountSig,
					CoinRef:   1,
				},
				"witness.outputs.[].token": {
					Format:    eip712.CSIGN_FIELD_FORMAT_RAW,
					Label:     "Tokens to receive",
					Signature: uniXTokenToReceiveSig,
				},
				"witness.outputs.[].amount": {
					Format:    eip712.CSIGN_FIELD_FORMAT_RAW,
					Label:     "Minimum amounts to receive",
					Signature: uniXMinAmountToReceiveSig,
				},
				"witness.outputs.[].recipient": {
					Format:    eip712.CSIGN_FIELD_FORMAT_RAW,
					Label:     "On Addresses",
					Signature: uniXOnAddressSig,
				},
				"deadline": {
					Format:    eip712.CSIGN_FIELD_FORMAT_DATETIME,
					Label:     "Approval expire",
					Signature: uniXApprovalExpire,
				},
			},
			ERC20Signatures: util.ERC20Sigs,
		},
	}
	eip712Sig, err := ethApp.SignEIP712Message(ctx, walletPath, message)

	if err != nil {
		logger.Error("unable to sign EIP712 message", "err", err)
		return
	}

	logger.Info("Signature", "R", log.HexDisplay(eip712Sig.R[:]), "S", log.HexDisplay(eip712Sig.S[:]), "V", eip712Sig.V)
}
