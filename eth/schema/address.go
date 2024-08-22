package schema

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ntchjb/ledger-go/adpu"
)

// Examples
// BIP32Path = "m'/44'/60'/2'/0/0"
// #1 Split string => ["m'", "44'", "60'", "2'", "0", "0"]
// #2 Remove "m" character => ["44'", "60'", "2'", "0", "0"]
// #3 Trim single quote => ["44", "60", "2", "0", "0"]
// #4 Add 1's bit at MSB position for those paths that have single quote
func SplitBIP32Paths(bip32Path string) ([]uint32, error) {
	var res []uint32
	paths := strings.Split(bip32Path, "/")
	// ignore leading 'm' character
	if paths[0] == "m'" {
		paths = paths[1:]
	}

	for _, path := range paths {
		numStr := strings.Trim(path, "'")
		num, err := strconv.ParseUint(numStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("unable to parse path num %s: %w", numStr, err)
		}

		if path[len(path)-1] == '\'' {
			num += 0x8000_0000
		}
		res = append(res, uint32(num))
	}

	return res, nil
}

type BIP32Path string

func (p *BIP32Path) Len() int {
	paths, err := SplitBIP32Paths(string(*p))
	if err != nil {
		return 0
	}

	return 1 + len(paths)*4
}

func (p *BIP32Path) MarshalADPU() ([]byte, error) {
	paths, err := SplitBIP32Paths(string(*p))

	if err != nil {
		return nil, fmt.Errorf("unable to split paths %s: %w", *p, err)
	}

	res := make([]byte, 1+len(paths)*4)

	res[0] = byte(len(paths))
	for i, num := range paths {
		binary.BigEndian.PutUint32(res[1+i*4:i*4+5], num)
	}

	return res, nil
}

type GetAddressRequest struct {
	// BIP-32 path for accessing an address in HD wallets
	BIP32Path BIP32Path
	ChainID   uint64
}

func (c *GetAddressRequest) MarshalADPU() ([]byte, error) {
	buf, err := adpu.Marshal(&c.BIP32Path)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal BIP-32 path: %w", err)
	}

	if c.ChainID > 0 {
		var chainIDBytes [8]byte
		binary.BigEndian.PutUint64(chainIDBytes[:], c.ChainID)
		buf = append(buf, chainIDBytes[:]...)
	}

	return buf, nil
}

type GetAddressResponse struct {
	// Public key of the address
	// Check the first byte to see format as follows
	// - 0x04: Uncompressed public key [0x04, (X value 32 bytes)..., (Y value 32 bytes)]
	// - 0x03: Compressed public key, where Y value is ODD [0x03, (X value 32 bytes)]
	// - 0x02: Compressed public key, where Y value is EVEN [0x02, (X value 32 bytes)]
	PublicKey PublicKey
	// Wallet address
	Address Address
	// Extension data generated by BIP-32 HD wallets for deriving child keys
	Chaincode ChainCode
}

func (g *GetAddressResponse) UnmarshalADPU(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data, cannot get public key length")
	}
	publicKeyLength := data[0]
	if len(data) <= 1+int(publicKeyLength) {
		return fmt.Errorf("data too short, cannot get address length")
	}
	addresslength := data[1+publicKeyLength]
	minDataLength := 2 + int(publicKeyLength) + int(addresslength)
	if len(data) < minDataLength {
		return fmt.Errorf("data too short, cannot get address length")
	}

	address := make([]byte, addresslength)
	copy(address, data[2+publicKeyLength:minDataLength])
	_, err := hex.Decode(g.Address[:], address)
	if err != nil {
		return fmt.Errorf("unable to decode address as hex: %w", err)
	}

	copy(g.PublicKey[:], data[1:1+publicKeyLength])
	if len(data) > minDataLength {
		copy(g.Chaincode[:], data[minDataLength:minDataLength+CHAIN_CODE_LENGTH])
	}

	return nil
}
