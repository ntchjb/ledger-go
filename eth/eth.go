package eth

import (
	"context"
	"fmt"

	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/eth/schema"
)

type EthereumApp interface {
	// Get Ledger Ethereum app's configurations
	GetConfiguration(ctx context.Context) (schema.Configuration, error)
	// Get address based on BIP-32 path string i.e. "m'/44'/60'/2'/0/0"
	GetAddress(ctx context.Context, bip32Path string, needHWConfirm bool, chaincode bool, chainID uint64) (schema.GetAddressResponse, error)
}

type ethereumAppImpl struct {
	proto adpu.Protocol
}

func NewEthereumApp(proto adpu.Protocol) EthereumApp {
	return &ethereumAppImpl{
		proto: proto,
	}
}

func (e *ethereumAppImpl) GetConfiguration(ctx context.Context) (schema.Configuration, error) {
	var conf schema.Configuration
	if err := adpu.Send(
		ctx,
		e.proto,
		ADPU_CLA, ADPU_INS_GET_CONFIGURATION, 0x00, 0x00,
		&adpu.EmptyData{},
		&conf,
	); err != nil {
		return conf, fmt.Errorf("unable to send get configuration to device: %w", err)
	}

	return conf, nil
}

func (e *ethereumAppImpl) GetAddress(ctx context.Context, bip32Path string, needHWConfirm bool, chaincode bool, chainID uint64) (schema.GetAddressResponse, error) {
	req := schema.GetAddressRequest{
		BIP32Path: bip32Path,
		ChainID:   chainID,
	}
	var res schema.GetAddressResponse

	var p1 uint8
	var p2 uint8
	if needHWConfirm {
		p1 = 0x01
	}
	if chaincode {
		p2 = 0x01
	}

	if err := adpu.Send(ctx, e.proto, ADPU_CLA, ADPU_INS_GET_PUBLIC_KEY, p1, p2, &req, &res); err != nil {
		return res, fmt.Errorf("unable to send get address to device: %w", err)
	}

	return res, nil
}
