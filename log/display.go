package log

import (
	"encoding/hex"
	"encoding/json"
)

type HexDisplay []byte

func (d HexDisplay) MarshalText() (text []byte, err error) {
	text = make([]byte, hex.EncodedLen(len(d)))
	hex.Encode(text, d)

	return
}

func (d HexDisplay) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(d))
}
