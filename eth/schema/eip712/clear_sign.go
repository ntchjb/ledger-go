package eip712

import (
	"errors"

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

type ClearSigning struct {
	Enabled bool

	ContractInfo CSignContract
	// Fields' key is path
	Fields          map[string]CSignField
	ERC20Signatures schema.ERC20Signatures
	// Key is CoinRef and Value is path
	// If Key is 255, it is domain's verifying contract
	CoinRefMap map[int]schema.Address
}

func (c *ClearSigning) ContractPayload() []byte {
	fieldCount := byte(len(c.Fields))
	return append(append(c.ContractInfo.DisplayPayload(), fieldCount), c.ContractInfo.SignaturePayload()...)
}
