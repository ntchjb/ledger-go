package eip712_test

import (
	"testing"

	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/stretchr/testify/assert"
)

func TestFieldDefinition_MarshalADPU(t *testing.T) {
	keyName := "field1"
	tests := []struct {
		name     string
		fieldDef eip712.FieldDefinition
		res      []byte
		err      error
	}{
		{
			name: "Success_AtomicType_FixedSize",
			fieldDef: eip712.FieldDefinition{
				TypeDescription: eip712.FieldTypeDescription{
					IsArray:         false,
					IsSizeSpecified: true,
					Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
				},
				TypeSize: 8, // uint64
				KeyName:  keyName,
			},
			res: []byte{
				0b0100_0010,
				0x08,
				0x06, 0x66, 0x69, 0x65, 0x6C, 0x64, 0x31,
			},
		},
		{
			name: "Success_AtomicType_DynamicSize",
			fieldDef: eip712.FieldDefinition{
				TypeDescription: eip712.FieldTypeDescription{
					IsArray:         false,
					IsSizeSpecified: false,
					Type:            eip712.FIELD_TYPE_DESC_TYPE_DYNAMIC_SIZED_BYTES,
				},
				KeyName: keyName,
			},
			res: []byte{
				0b0000_0111,
				0x06, 0x66, 0x69, 0x65, 0x6C, 0x64, 0x31,
			},
		},
		{
			name: "Success_AtomicType_DynamicSize_String",
			fieldDef: eip712.FieldDefinition{
				TypeDescription: eip712.FieldTypeDescription{
					IsArray:         false,
					IsSizeSpecified: false,
					Type:            eip712.FIELD_TYPE_DESC_TYPE_STRING,
				},
				KeyName: keyName,
			},
			res: []byte{
				0b0000_0101,
				0x06, 0x66, 0x69, 0x65, 0x6C, 0x64, 0x31,
			},
		},
		{
			name: "Success_AtomicType_Array_uint8[3][]",
			fieldDef: eip712.FieldDefinition{
				TypeDescription: eip712.FieldTypeDescription{
					IsArray:         true,
					IsSizeSpecified: true,
					Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
				},
				TypeSize: 1,
				ArrayLevels: []eip712.FieldArrayLevel{
					{
						Type:           eip712.STRUCT_DEF_ARRAY_TYPE_FIXED,
						FixedArraySize: 3,
					},
					{
						Type: eip712.STRUCT_DEF_ARRAY_TYPE_DYNAMIC,
					},
				},
				KeyName: keyName,
			},
			res: []byte{
				0b1100_0010,
				0x01,
				0x02,
				0x01, 0x03,
				0x00,
				0x06, 0x66, 0x69, 0x65, 0x6C, 0x64, 0x31,
			},
		},
		{
			name: "Success_CustomType",
			fieldDef: eip712.FieldDefinition{
				TypeDescription: eip712.FieldTypeDescription{
					IsArray:         false,
					IsSizeSpecified: false,
					Type:            eip712.FIELD_TYPE_DESC_TYPE_CUSTOM,
				},
				TypeName: "Mail",
				KeyName:  keyName,
			},
			res: []byte{
				0b0000_0000,
				0x04, 0x4D, 0x61, 0x69, 0x6C,
				0x06, 0x66, 0x69, 0x65, 0x6C, 0x64, 0x31,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := test.fieldDef.MarshalADPU()

			assert.ErrorIs(t, err, test.err)
			assert.Equal(t, test.res, res)
		})
	}
}
