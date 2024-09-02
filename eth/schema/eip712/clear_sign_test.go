package eip712_test

import (
	"testing"

	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/stretchr/testify/assert"
)

func TestCSignContract(t *testing.T) {
	contract := eip712.CSignContract{
		Label:     "contractName",
		Signature: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
	}

	assert.Equal(t, append([]byte{0x0C}, []byte(contract.Label)...), contract.DisplayPayload())
	assert.Equal(t, []byte{0x06, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05}, contract.SignaturePayload())
}

func TestCSignField_Payload(t *testing.T) {
	type fields struct {
		format    eip712.CSignFieldFormat
		label     string
		signature []byte
		coinRef   int
	}

	type args struct {
		registeredCoinRef map[int]uint8
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		resp   []byte
		err    error
	}{
		{
			name: "Success_TokenField",
			fields: fields{
				format:    "token",
				label:     "DAI",
				signature: []byte{0x01, 0x02, 0x03},
				coinRef:   0,
			},
			args: args{
				registeredCoinRef: map[int]uint8{
					0: 1,
				},
			},
			resp: []byte{
				0x01,
				0x03, 0x01, 0x02, 0x03,
			},
		},
		{
			name: "Success_RawField",
			fields: fields{
				format:    "raw",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
			},
			args: args{
				registeredCoinRef: map[int]uint8{
					0: 1,
				},
			},
			resp: append(append(
				[]byte{0x0C}, []byte("AwesomeLabel")...),
				0x03,
				0x01, 0x02, 0x03,
			),
			err: nil,
		},
		{
			name: "Success_DateTime",
			fields: fields{
				format:    "datetime",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
			},
			args: args{
				registeredCoinRef: map[int]uint8{
					0: 1,
				},
			},
			resp: append(
				append([]byte{0x0C}, []byte("AwesomeLabel")...),
				0x03,
				0x01, 0x02, 0x03,
			),
			err: nil,
		},
		{
			name: "Success_Amount",
			fields: fields{
				format:    "amount",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
				coinRef:   1,
			},
			args: args{
				registeredCoinRef: map[int]uint8{
					1: 3,
				},
			},
			resp: append(
				append([]byte{0x0C}, []byte("AwesomeLabel")...),
				0x03,
				0x03, 0x01, 0x02, 0x03,
			),
			err: nil,
		},
		{
			name: "Error_TokenField_RegisteredCoinNotFound",
			fields: fields{
				format:    "amount",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
				coinRef:   1,
			},
			args: args{
				registeredCoinRef: map[int]uint8{},
			},
			resp: nil,
			err:  eip712.ErrRegisteredCoinRefNotFound,
		},
		{
			name: "Error_AmountField_RegisteredCoinNotFound",
			fields: fields{
				format:    "token",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
				coinRef:   1,
			},
			args: args{
				registeredCoinRef: nil,
			},
			resp: nil,
			err:  eip712.ErrRegisteredCoinRefNotFound,
		},
		{
			name: "Error_UnknownField",
			fields: fields{
				format:    "unknown",
				label:     "AwesomeLabel",
				signature: []byte{0x01, 0x02, 0x03},
				coinRef:   1,
			},
			args: args{
				registeredCoinRef: nil,
			},
			resp: nil,
			err:  eip712.ErrUnknownFieldFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cSignField := eip712.CSignField{
				Format:    test.fields.format,
				Label:     test.fields.label,
				Signature: test.fields.signature,
				CoinRef:   test.fields.coinRef,
			}

			res, err := cSignField.Payload(test.args.registeredCoinRef)

			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.resp, res)
		})
	}
}

func TestClearSigning_ContractPayload(t *testing.T) {
	cs := eip712.ClearSigning{
		Enabled: true,
		ContractInfo: eip712.CSignContract{
			Label:     "OrderManager",
			Signature: []byte{0x01, 0x02, 0x03},
		},
		Fields: map[string]eip712.CSignField{
			"asset":  {},
			"asset2": {},
		},
	}

	cPayload := cs.ContractPayload()

	assert.Equal(t,
		cPayload,
		append(
			append([]byte{0x0C}, []byte("OrderManager")...),
			2,
			0x03, 0x01, 0x02, 0x03,
		),
	)
}

func TestCSignField_Action(t *testing.T) {
	tests := []struct {
		name   string
		format eip712.CSignFieldFormat
		action eip712.Action
		err    error
	}{
		{
			name:   "Raw format",
			format: eip712.CSIGN_FIELD_FORMAT_RAW,
			action: eip712.ACTION_RAW,
		},
		{
			name:   "Datetime format",
			format: eip712.CSIGN_FIELD_FORMAT_DATETIME,
			action: eip712.ACTION_DATETIME,
		},
		{
			name:   "Token format",
			format: eip712.CSIGN_FIELD_FORMAT_TOKEN,
			action: eip712.ACTION_AMOUNT_TOKEN_JOIN,
		},
		{
			name:   "Amount format",
			format: eip712.CSIGN_FIELD_FORMAT_AMOUNT,
			action: eip712.ACTION_AMOUNT_VALUE_JOIN,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			csf := eip712.CSignField{
				Format: test.format,
			}

			action, err := csf.Action()

			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.action, action)
		})
	}
}
