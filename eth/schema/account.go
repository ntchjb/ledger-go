package schema

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

const (
	ADDRESS_LENGTH    int = 20
	PUBLIC_KEY_LENGTH int = 65
	CHAIN_CODE_LENGTH int = 32
)

type PublicKey [PUBLIC_KEY_LENGTH]byte

func (p *PublicKey) String() string {
	return "0x" + hex.EncodeToString(p[:])
}

type Address [ADDRESS_LENGTH]byte

func (a *Address) hexBytes() [ADDRESS_LENGTH*2 + 2]byte {
	var res [ADDRESS_LENGTH*2 + 2]byte
	copy(res[:2], []byte("0x"))
	hex.Encode(res[2:], a[:])

	return res
}

func (a *Address) String() string {
	hexBytes := a.hexBytes()
	hasher := sha3.NewLegacyKeccak256()

	_, err := hasher.Write(hexBytes[2:])
	if err != nil {
		return ""
	}
	hash := hasher.Sum(nil)

	for i := 2; i < len(hexBytes); i++ {
		hashNibble := hash[(i-2)/2] >> (((i & 1) ^ 1) << 2)
		if hexBytes[i] > '9' && (hashNibble&0x0F > 7) {
			hexBytes[i] -= 32
		}
	}

	return string(hexBytes[:])
}

type ChainCode [CHAIN_CODE_LENGTH]byte

func (c *ChainCode) String() string {
	return "0x" + hex.EncodeToString(c[:])
}
