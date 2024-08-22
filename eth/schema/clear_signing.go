package schema

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
	ExternalPlugin ExternalPluginResolution
	// Plugin payloads for displaying contract calls information
	Plugin []PluginResolution
	// Show domain address information on Ledger display i.e. ENS domain name
	Domains []DomainResolution
}
