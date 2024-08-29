package eip712_test

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/stretchr/testify/assert"
)

func TestMessage_SetCoinRefMap(t *testing.T) {
	type fields struct {
		chainID            uint64
		clearSigningFields map[string]eip712.CSignField
	}
	type args struct {
		signingData eip712.StructItem
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		err        error
		coinRefMap map[int]schema.Address
	}{
		{
			name: "Success",
			fields: fields{
				chainID: 10,
				clearSigningFields: map[string]eip712.CSignField{
					"from": {
						Format:  eip712.CSIGN_FIELD_FORMAT_TOKEN,
						CoinRef: 0,
					},
					"to": {
						Format:  eip712.CSIGN_FIELD_FORMAT_TOKEN,
						CoinRef: 2,
					},
					"nested.createdBy": {
						Format:  eip712.CSIGN_FIELD_FORMAT_TOKEN,
						CoinRef: 1,
					},
					"unreachable": {
						Format:  eip712.CSIGN_FIELD_FORMAT_TOKEN,
						CoinRef: 10,
					},
					"nested.presentedBy": {
						Format:  eip712.CSIGN_FIELD_FORMAT_AMOUNT,
						CoinRef: 255,
					},
				},
			},
			args: args{
				signingData: eip712.StructItem{
					TypeName: "Order",
					Members: []eip712.StructItemMember{
						{
							Name: "from",
							Item: eip712.AtomicItem{
								Item: eip712.AddressData{
									0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
									0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
								},
							},
						},
						{
							Name: "to",
							Item: eip712.AtomicItem{
								Item: eip712.AddressData{
									0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
									0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
								},
							},
						},
						{
							Name: "ayo",
							Item: eip712.AtomicItem{
								Item: eip712.AddressData{
									0xEE, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
									0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
								},
							},
						},
						{
							Name: "nested",
							Item: eip712.StructItem{
								TypeName: "Creator",
								Members: []eip712.StructItemMember{
									{
										Name: "createdBy",
										Item: eip712.AtomicItem{
											Item: eip712.AddressData{
												0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
												0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a,
											},
										},
									},
									{
										Name: "presentedBy",
										Item: eip712.AtomicItem{
											Item: eip712.NumberData{
												Num:     uint256.NewInt(10),
												NumBits: 8,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			err: nil,
			coinRefMap: map[int]schema.Address{
				0: {
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
				},
				1: {
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
					0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a,
				},
				2: {
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a,
				},
				255: {
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			message := eip712.Message{
				Domain: eip712.Domain{
					VerifyingContract: schema.Address{
						0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
						0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					},
					ChainID: uint256.NewInt(test.fields.chainID),
				},
				ClearSigning: eip712.ClearSigning{
					Fields: test.fields.clearSigningFields,
				},
			}
			err := message.SetCoinRefMap(test.args.signingData)

			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.coinRefMap, message.ClearSigning.CoinRefMap)
		})
	}
}
