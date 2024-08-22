package schema

import (
	"fmt"
)

type GetConfigurationResponse struct {
	ArbitraryDataEnabled       bool
	ERC20ProvisioningNecessary bool
	StarkEnabled               bool
	StarkV2Supported           bool
	Version                    string
}

func (c *GetConfigurationResponse) UnmarshalADPU(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("data too short, expected %d, got %d", 4, len(data))
	}

	c.ArbitraryDataEnabled = data[0]&0x01 == 0x01
	c.ERC20ProvisioningNecessary = data[0]&0x02 == 0x02
	c.StarkEnabled = data[0]&0x04 == 0x04
	c.StarkV2Supported = data[0]&0x08 == 0x08
	c.Version = fmt.Sprintf("%d.%d.%d", data[1], data[2], data[3])

	return nil
}
