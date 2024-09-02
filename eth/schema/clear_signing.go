package schema

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/holiman/uint256"
)

type Challenge [4]byte

func (c *Challenge) UnmarshalADPU(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("data is too short, expected 4, got %d", len(data))
	}

	copy((*c)[:], data)

	return nil
}

type DomainNameBlob []byte

func (d *DomainNameBlob) MarshalADPU() ([]byte, error) {
	var res []byte
	var length [2]byte

	binary.BigEndian.PutUint16(length[:], uint16(len(*d)))
	res = append(length[:], (*d)...)

	return res, nil
}

type ExternalPluginResolution struct {
	Payload   []byte
	Signature []byte
}

type SupportedRegistry string

const (
	DOMAIN_REGISTRY_ENS SupportedRegistry = "ens"
)

type DomainType uint8

const (
	DOMAIN_TYPE_FORWARD DomainType = iota
	DOMAIN_TYPE_REVERSED
)

type DomainResolution struct {
	Registry SupportedRegistry
	Domain   string
	Address  string
	Type     DomainType
}

type DomainSignatureResponse struct {
	Payload string `json:"payload"`
}

type ERC20TokenResolution []byte
type PluginResolution []byte
type NFTResolution []byte

type ClearSigningResolution struct {
	// ERC20 tokens payloads for displaying ERC20 transaction info
	ERC20Tokens []ERC20TokenResolution
	// NFTs payloads for displaying NFT transaction info
	NFTs []NFTResolution
	// External plugin payloads for displaying contract calls information
	// using external plugins
	ExternalPlugin []ExternalPluginResolution
	// Plugin payloads for displaying contract calls information
	Plugin []PluginResolution
	// Show domain address information on Ledger display i.e. ENS domain name
	Domains []DomainResolution
}

type ProvideERC20InfoResponse byte

func (r *ProvideERC20InfoResponse) UnmarshalADPU(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data is too short, expected 1, got 0")
	}

	*r = ProvideERC20InfoResponse(data[0])

	return nil
}

type CSignTokenInfo struct {
	ContractAddress Address
	Ticker          string
	Decimals        uint32
	ChainID         uint32
	Signature       []byte

	Raw []byte
}

type ERC20Signatures map[[24]byte]CSignTokenInfo

func (e ERC20Signatures) FindByChainIDAndAddress(chainID *uint256.Int, address Address) (CSignTokenInfo, bool) {
	var key [24]byte
	bChainID := chainID.Bytes32()
	copy(key[:4], bChainID[len(bChainID)-4:])
	copy(key[4:], address[:])

	res, ok := e[key]
	return res, ok
}

func ParseERC20SignatureBlobs(blob string) (ERC20Signatures, error) {
	res := make(map[[24]byte]CSignTokenInfo)
	blobBytes, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return res, fmt.Errorf("unable to decode base64 blob: %w", err)
	}
	for i := 0; i < len(blobBytes); {
		var tokenInfo CSignTokenInfo
		var key [24]byte

		recordLength := binary.BigEndian.Uint32(blobBytes[i : i+4])
		i += 4

		record := blobBytes[i : i+int(recordLength)]

		recIdx := 0
		tickerLength := int(record[recIdx])
		recIdx += 1

		tokenInfo.Ticker = string(record[recIdx : recIdx+tickerLength])
		recIdx += tickerLength

		copy(tokenInfo.ContractAddress[:], record[recIdx:recIdx+20])
		copy(key[4:], record[recIdx:recIdx+20])
		recIdx += 20

		tokenInfo.Decimals = binary.BigEndian.Uint32(record[recIdx : recIdx+4])
		recIdx += 4

		tokenInfo.ChainID = binary.BigEndian.Uint32(record[recIdx : recIdx+4])
		copy(key[:4], record[recIdx:recIdx+4])
		recIdx += 4

		tokenInfo.Signature = record[recIdx:]
		tokenInfo.Raw = record

		res[key] = tokenInfo

		i += int(recordLength)
	}

	return res, nil
}
