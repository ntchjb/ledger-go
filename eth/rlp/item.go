package rlp

import (
	"encoding/binary"
	"errors"

	"github.com/holiman/uint256"
)

type ItemType uint8

const (
	ITEM_TYPE_STRING ItemType = iota
	ITEM_TYPE_LIST
)

var (
	ErrEmptyInput       = errors.New("empty input")
	ErrInvalidRLPFormat = errors.New("invalid RLP format")
)

type ItemInfo struct {
	// Start position of item's content
	Offset int
	// Length of item's content
	Length int
	// Item type
	Type ItemType
}

type Item struct {
	Data []byte
	List []Item
}

func (t *Item) String() string {
	return string(t.Data)
}

func (t *Item) Uint64() uint64 {
	size := 8
	var numBytes [8]byte

	copy(numBytes[max(0, size-len(t.Data)):], t.Data[max(0, len(t.Data)-8):len(t.Data)])

	return binary.BigEndian.Uint64(numBytes[:])
}

func (t *Item) Uint256() *uint256.Int {
	return new(uint256.Int).SetBytes(t.Data)
}

func (t *Item) getLengthOfDataLengthNumber(dataLength int) int {
	if dataLength <= 0xFF {
		return 1
	} else if dataLength <= 0xFFFF {
		return 2
	} else if dataLength <= 0xFFFFFF {
		return 3
	} else if dataLength <= 0xFFFFFFFF {
		return 4
	} else if dataLength <= 0xFFFFFFFFFF {
		return 5
	} else if dataLength <= 0xFFFFFFFFFFFF {
		return 6
	} else if dataLength <= 0xFFFFFFFFFFFFFF {
		return 7
	} else {
		return 8
	}
}

func (t *Item) Len() int {
	if t.Data != nil {
		if len(t.Data) == 0 {
			return 1
		} else if len(t.Data) == 1 {
			if t.Data[0] <= 0x7F {
				// An ASCII character
				return 1
			} else {
				// A short string with length 1
				return 2
			}
		} else if len(t.Data) <= 55 {
			// A short string
			return 1 + len(t.Data)
		} else if len(t.Data) > 55 {
			// A long string
			return 1 + t.getLengthOfDataLengthNumber(len(t.Data)) + len(t.Data)
		}
	} else if t.List != nil {
		var listLength int
		for _, item := range t.List {
			listLength += item.Len()
		}
		if listLength <= 55 {
			return 1 + listLength
		} else {
			return 1 + t.getLengthOfDataLengthNumber(listLength) + listLength
		}
	}

	return 0
}
