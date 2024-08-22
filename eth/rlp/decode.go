package rlp

import (
	"encoding/binary"
	"fmt"
)

func Decode(encoded []byte) (Item, int, error) {
	var res Item

	info, err := decodeItem(encoded)
	if err != nil {
		return Item{}, 0, err
	}

	if info.Type == ITEM_TYPE_STRING {
		res.Data = encoded[info.Offset : info.Offset+info.Length]
	} else if info.Type == ITEM_TYPE_LIST {
		res.List = []Item{}
		encoded = encoded[info.Offset:]
		totalLength := 0
		for totalLength < info.Length {
			item, length, err := Decode(encoded)
			if err != nil {
				return Item{}, 0, fmt.Errorf("unable to decode list, offset: %d, length: %d, err: %w", info.Offset, info.Length, err)
			}

			res.List = append(res.List, item)

			totalLength += length
			encoded = encoded[length:]
		}

	}

	return res, info.Offset + info.Length, nil
}

func getDynamicDataLength(lengthBytes []byte) uint64 {
	var uint64Bytes [8]byte
	copy(uint64Bytes[8-len(lengthBytes):], lengthBytes)

	return binary.BigEndian.Uint64(uint64Bytes[:])
}

func decodeItem(encoded []byte) (ItemInfo, error) {
	var res ItemInfo
	if len(encoded) == 0 {
		return res, ErrEmptyInput
	}
	prefix := encoded[0]

	if prefix <= 0x7F {
		// An ASCII character, the first byte represents data itself
		// [(data, 1 byte)]
		res.Offset = 0
		res.Length = 1
		res.Type = ITEM_TYPE_STRING
	} else if prefix <= 0xB7 && len(encoded) > int(prefix-0x80) {
		// A short data with length 0-55 bytes, first byte is data length, followed by data as bytes
		// [(data length, 1 bytes), (data)...]
		res.Offset = 1
		res.Length = int(prefix - 0x80)
		res.Type = ITEM_TYPE_STRING
	} else if prefix <= 0xBF && len(encoded) > int(prefix-0xB7) && uint64(len(encoded)) > uint64(prefix-0xB7)+getDynamicDataLength(encoded[1:prefix-0xB7]) {
		// A long data,
		// [(length of data length number, 1 byte), (data length, 0-8 bytes), (data)...]
		// first byte tell number of bytes to store string length number,
		// followed by data length number, big-endian encoded (maximum of 8 bytes, or unsigned int 64-bit),
		// followed by data as bytes
		res.Offset = 1 + int(prefix-0xB7)
		res.Length = int(getDynamicDataLength(encoded[1 : 1+prefix-0xB7]))
		res.Type = ITEM_TYPE_STRING
	} else if prefix <= 0xF7 && len(encoded) > int(prefix-0xC0) {
		// A short list with data length between 0-55 bytes,
		// [(data length, 1 bytes), (data)...]
		// each item in the list is RLP-encoded and concatenated together
		res.Offset = 1
		res.Length = int(prefix - 0xC0)
		res.Type = ITEM_TYPE_LIST
	} else if len(encoded) > int(prefix-0xF7) && uint64(len(encoded)) > uint64(prefix-0xF7)+getDynamicDataLength(encoded[1:prefix-0xF7]) {
		// A long list,
		// [(length of data length number, 1 byte), (data length, 0-8 bytes), (data)...]
		// the first byte indicates number of bytes to store list data length number
		// followed by data length number, big-endian encoded (maximum of 8 bytes, or unsigned int 64-bit),
		// followed by data as bytes
		// each item in the list is RLP-encoded and concatenated together
		res.Offset = 1 + int(prefix-0xF7)
		res.Length = int(getDynamicDataLength(encoded[1 : 1+prefix-0xF7]))
		res.Type = ITEM_TYPE_LIST
	} else {
		return res, fmt.Errorf("%w, got prefix 0x%x, remaining length: %d", ErrInvalidRLPFormat, prefix, len(encoded))
	}

	return res, nil
}
