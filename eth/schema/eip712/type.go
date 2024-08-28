package eip712

import "fmt"

type Component uint8

const (
	// Struct name
	TYPE_COMPONENT_NAME Component = 0x00
	// Field name
	TYPE_COMPONENT_FIELD Component = 0xFF

	DATA_COMPONENT_ROOT   Component = 0x00
	DATA_COMPONENT_ARRAY  Component = 0x0F
	DATA_COMPONENT_ATOMIC Component = 0xFF
)

type Action uint8

const (
	ACTION_ACTIVATE          Action = 0x00
	ACTION_MESSAGE_INFO      Action = 0x0F
	ACTION_DATETIME          Action = 0xFC
	ACTION_AMOUNT_TOKEN_JOIN Action = 0xFD
	ACTION_AMOUNT_VALUE_JOIN Action = 0xFE
	ACTION_RAW               Action = 0xFF
)

type StructDefArrayType uint8

const (
	STRUCT_DEF_ARRAY_TYPE_DYNAMIC StructDefArrayType = 0x00
	STRUCT_DEF_ARRAY_TYPE_FIXED   StructDefArrayType = 0x01
)

const (
	DOMAIN_TYPE_NAME = "EIP712Domain"
)

type FieldArrayLevel struct {
	Type           StructDefArrayType
	FixedArraySize uint8
}

type FieldType byte

const (
	FIELD_TYPE_DESC_IS_ARRAY          FieldType = 0b1000_0000
	FIELD_TYPE_DESC_IS_SIZE_SPECIFIED FieldType = 0b0100_0000

	FIELD_TYPE_DESC_TYPE_CUSTOM              FieldType = 0x00
	FIELD_TYPE_DESC_TYPE_INT                 FieldType = 0x01
	FIELD_TYPE_DESC_TYPE_UINT                FieldType = 0x02
	FIELD_TYPE_DESC_TYPE_ADDRESS             FieldType = 0x03
	FIELD_TYPE_DESC_TYPE_BOOL                FieldType = 0x04
	FIELD_TYPE_DESC_TYPE_STRING              FieldType = 0x05
	FIELD_TYPE_DESC_TYPE_FIXED_SIZE_BYTES    FieldType = 0x06
	FIELD_TYPE_DESC_TYPE_DYNAMIC_SIZED_BYTES FieldType = 0x07
)

type FieldTypeDescription struct {
	// Is this struct type an array?
	IsArray bool
	// Is this type has indicated size? i.e. uint8, int256, bytes32
	IsSizeSpecified bool
	// Main type of this struct i.e. `int[2]` => main type is `int`
	Type FieldType
}

// A member of EIP-712 `encodeType`
type FieldDefinition struct {
	// Type description of given field
	TypeDescription FieldTypeDescription
	// Type name, for custom type
	TypeName string
	// Type size in bytes (1-255 bytes) for types with fixed bits size specified
	// i.e. uint256, int32, bytes20
	TypeSize uint8
	// Type and length of the array in every dimensions
	// i.e. int[3][][4]
	// Item #0: Fixed array with size 3
	// Item #1: Dynamic array
	// Item #2: Fixed array with size 4
	ArrayLevels []FieldArrayLevel
	// Field name
	KeyName string
}

func (s *FieldDefinition) MarshalADPU() ([]byte, error) {
	var res []byte

	var typeDesc byte
	// #1: Set type description (1 byte)
	if s.TypeDescription.IsArray {
		typeDesc |= byte(FIELD_TYPE_DESC_IS_ARRAY)
	}
	if s.TypeDescription.IsSizeSpecified {
		typeDesc |= byte(FIELD_TYPE_DESC_IS_SIZE_SPECIFIED)
	}
	typeDesc |= byte(s.TypeDescription.Type)
	res = append(res, typeDesc)

	// #2: Set type name length, and type name string for custom type (1 + len(s.TypeName) bytes)
	if s.TypeDescription.Type == FIELD_TYPE_DESC_TYPE_CUSTOM {
		buf := make([]byte, 1+len(s.TypeName))
		if len(s.TypeName) > 255 {
			return nil, fmt.Errorf("type name is too long, expected <256, got %d", len(s.TypeName))
		}
		buf[0] = byte(len(s.TypeName))
		copy(buf[1:], []byte(s.TypeName))
		res = append(res, buf...)
	}

	// #3: Set type size if type size specified in type name
	if s.TypeDescription.IsSizeSpecified {
		res = append(res, s.TypeSize)
	}

	// #4: Set array type and its size for every dimensions, if this is array
	if s.TypeDescription.IsArray {
		if len(s.ArrayLevels) > 255 {
			return nil, fmt.Errorf("array levels is too deep, expected <256, got %d", len(s.ArrayLevels))
		}
		res = append(res, byte(len(s.ArrayLevels)))
		for _, level := range s.ArrayLevels {
			res = append(res, byte(level.Type))
			if level.Type == STRUCT_DEF_ARRAY_TYPE_FIXED {
				res = append(res, level.FixedArraySize)
			}
		}
	}

	// #5: Set field name
	if len(s.KeyName) > 255 {
		return nil, fmt.Errorf("key name is too long, expected <256, got %d", len(s.KeyName))
	}
	res = append(res, byte(len(s.KeyName)))
	res = append(res, []byte(s.KeyName)...)

	return res, nil
}
