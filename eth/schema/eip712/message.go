package eip712

import (
	"fmt"

	"github.com/ntchjb/ledger-go/eth/schema"
)

type TypeStruct struct {
	// Type name
	Name string
	// Type members
	Members []FieldDefinition
}

type TypeStructs []TypeStruct

type Message struct {
	Types   TypeStructs
	Domain  Domain
	Primary StructItem

	ClearSigning ClearSigning
}

func (m *Message) SetCoinRefMap(signingData StructItem) error {
	m.ClearSigning.CoinRefMap = make(map[int]schema.Address)

	if err := signingData.Walk("", func(path string, item Item) error {
		fieldInfo, ok := m.ClearSigning.Fields[path]
		if !ok {
			return nil
		}
		if item.Type() != DATA_COMPONENT_ATOMIC {
			return fmt.Errorf("unsupported token type, supported atomic, got %d", item.Type())
		}

		if fieldInfo.Format == CSIGN_FIELD_FORMAT_TOKEN {
			token, ok := item.(AtomicItem).Item.(AddressData)
			if !ok {
				return fmt.Errorf("token is not address, cannot convert to AddressData")
			}
			m.ClearSigning.CoinRefMap[fieldInfo.CoinRef] = schema.Address(token)
		} else if fieldInfo.Format == CSIGN_FIELD_FORMAT_AMOUNT && fieldInfo.CoinRef == 255 {
			m.ClearSigning.CoinRefMap[fieldInfo.CoinRef] = m.Domain.VerifyingContract
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error occurred during setting coin ref map: %w", err)
	}

	return nil
}
