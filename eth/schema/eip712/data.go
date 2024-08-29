package eip712

import (
	"encoding/binary"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/ntchjb/ledger-go/eth/schema"
)

type WalkReader func(path string, item Item) error

type Item interface {
	Type() Component
	Walk(path string, fn WalkReader) error
	DataCommand() DataCommand
}

type Domain struct {
	Name              string
	Version           string
	ChainID           *uint256.Int
	VerifyingContract schema.Address
	Salt              [32]byte
}

func (d *Domain) IsSaltExists() bool {
	isSaltExist := false
	for _, num := range d.Salt {
		if num != 0 {
			isSaltExist = true
			break
		}
	}

	return isSaltExist
}

func (d *Domain) TypeStruct() TypeStruct {
	res := TypeStruct{
		Name: DOMAIN_TYPE_NAME,
		Members: []FieldDefinition{
			{
				TypeDescription: FieldTypeDescription{
					Type: FIELD_TYPE_DESC_TYPE_STRING,
				},
				KeyName: "name",
			},
			{
				TypeDescription: FieldTypeDescription{
					Type: FIELD_TYPE_DESC_TYPE_STRING,
				},
				KeyName: "version",
			},
			{
				TypeDescription: FieldTypeDescription{
					IsSizeSpecified: true,
					Type:            FIELD_TYPE_DESC_TYPE_UINT,
				},
				TypeSize: 32,
				KeyName:  "chainId",
			},
			{
				TypeDescription: FieldTypeDescription{
					Type: FIELD_TYPE_DESC_TYPE_ADDRESS,
				},
				KeyName: "verifyingContract",
			},
		},
	}

	if d.IsSaltExists() {
		res.Members = append(res.Members, FieldDefinition{
			TypeDescription: FieldTypeDescription{
				IsSizeSpecified: true,
				Type:            FIELD_TYPE_DESC_TYPE_FIXED_SIZE_BYTES,
			},
			TypeSize: 32,
			KeyName:  "salt",
		})
	}

	return res
}

func (d *Domain) StructItem() StructItem {
	var res StructItem
	res.TypeName = DOMAIN_TYPE_NAME
	res.Members = []StructItemMember{
		{
			Name: "name",
			Item: AtomicItem{
				Item: StringData(d.Name),
			},
		},
		{
			Name: "version",
			Item: AtomicItem{
				Item: StringData(d.Version),
			},
		},
		{
			Name: "chainId",
			Item: AtomicItem{
				Item: NumberData{
					Num:     d.ChainID,
					NumBits: 256,
					Signed:  false,
				},
			},
		},
		{
			Name: "verifyingContract",
			Item: AtomicItem{
				Item: AddressData(d.VerifyingContract),
			},
		},
	}

	if d.IsSaltExists() {
		res.Members = append(res.Members, StructItemMember{
			Name: "salt",
			Item: AtomicItem{
				Item: BytesData{
					FixedSize: 32,
					Data:      d.Salt[:],
				},
			},
		})
	}

	return res
}

type StructItemMember struct {
	Name string
	Item Item
}

type StructItem struct {
	TypeName string
	Members  []StructItemMember
}

func (c StructItem) DataCommand() DataCommand {
	return DataCommand{
		Component: DATA_COMPONENT_ROOT,
		Value:     []byte(c.TypeName),
	}
}

func (c StructItem) Type() Component {
	return DATA_COMPONENT_ROOT
}

func (c StructItem) Walk(path string, fn WalkReader) error {
	if err := fn(path, c); err != nil {
		return fmt.Errorf("callback failed, path: %s, name: %s, itemType: %d, err: %w", path, c.TypeName, c.Type(), err)
	}

	for _, member := range c.Members {
		newPath := path
		if path != "" {
			newPath += "."
		}
		if err := member.Item.Walk(newPath+member.Name, fn); err != nil {
			return fmt.Errorf("callback failed, path: %s, name: %s, itemType: %d, err: %w", newPath+member.Name, member.Name, member.Item.Type(), err)
		}

	}

	return nil
}

type ArrayItem []Item

func (c ArrayItem) DataCommand() DataCommand {
	return DataCommand{
		Component: DATA_COMPONENT_ARRAY,
		Value:     []byte{byte(len(c))},
	}
}

func (c ArrayItem) Type() Component {
	return DATA_COMPONENT_ARRAY
}

func (c ArrayItem) Walk(path string, fn WalkReader) error {
	if err := fn(path, c); err != nil {
		return fmt.Errorf("callback failed, path: %s, name: (Array), itemType: %d, err: %w", path, c.Type(), err)
	}

	for _, item := range c {
		newPath := path
		if path != "" {
			newPath += "."
		}
		if err := item.Walk(newPath+"[]", fn); err != nil {
			return fmt.Errorf("callback failed, path: %s, name: (Array Item), itemType: %d, err: %w", newPath+"[]", item.Type(), err)
		}
	}

	return nil
}

type AtomicEncoder interface {
	Encode() []byte
}

type AtomicItem struct {
	Item AtomicEncoder
}

func (c AtomicItem) DataCommand() DataCommand {
	cmd := DataCommand{
		Component: DATA_COMPONENT_ATOMIC,
	}
	data := c.Item.Encode()

	if len(data) == 0 {
		return cmd
	}

	res := make([]byte, len(data)+2)
	binary.BigEndian.PutUint16(res[:2], uint16(len(data)))
	copy(res[2:], data)
	cmd.Value = res

	return cmd
}

func (c AtomicItem) Type() Component {
	return DATA_COMPONENT_ATOMIC
}

func (c AtomicItem) Walk(path string, fn WalkReader) error {
	if err := fn(path, c); err != nil {
		return fmt.Errorf("callback failed, path: %s, name: (Atomic Item), itemType: %d, err: %w", path, c.Type(), err)
	}

	return nil
}

// NumberData is a field representing uint8,...,uint256, int8,...,int256
type NumberData struct {
	Num *uint256.Int
	// Can this number be negative (int)?
	NumBits uint16
	Signed  bool
}

func (f NumberData) Encode() []byte {
	numBytes := f.NumBits / 8
	if numBytes > 32 {
		numBytes = 32
	}

	num := f.Num
	if f.Signed {
		num = new(uint256.Int).Not(f.Num)
		num = num.Add(num, uint256.NewInt(1))
	}
	b := num.Bytes32()
	return b[32-numBytes:]
}

type BoolData bool

func (f BoolData) Encode() []byte {
	if f {
		return []byte{0x01}
	}

	return []byte{0x00}
}

type AddressData schema.Address

func (f AddressData) Encode() []byte {
	return f[:]
}

type StringData string

func (f StringData) Encode() []byte {
	return []byte(f)
}

type BytesData struct {
	// Number of bytes for fixed bytes type i.e. `bytes20`, `bytes32`, etc.
	// 0 if it's dynamic `bytes` type
	FixedSize uint8
	Data      []byte
}

func (f BytesData) Encode() []byte {
	if f.FixedSize > 0 {
		return f.Data[:f.FixedSize]
	}

	return f.Data
}
