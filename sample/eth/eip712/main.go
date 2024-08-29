package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/holiman/uint256"
	"github.com/ntchjb/gohid/hid"
	"github.com/ntchjb/gohid/manager"
	"github.com/ntchjb/gohid/usb"
	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/device"
	"github.com/ntchjb/ledger-go/eth"
	"github.com/ntchjb/ledger-go/eth/schema"
	"github.com/ntchjb/ledger-go/eth/schema/eip712"
	"github.com/ntchjb/ledger-go/log"
)

func main() {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	usbCtx := usb.NewGOUSBContext()
	man := manager.NewDeviceManager(usbCtx, logger)

	defer func() {
		if err := man.Close(); err != nil {
			logger.Error("unable to close device manager", "err", err)
		}
	}()

	// Open Ledger Nano S
	// hidDevice, err := man.Open(0x2C97, 0x1011, hid.DeviceConfig{
	// 	StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	// })

	// Open Ledger Nano S Plus
	hidDevice, err := man.Open(0x2C97, 0x5011, hid.DeviceConfig{
		StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	})

	if err != nil {
		logger.Error("unable to open device", "err", err)
		return
	}
	defer hidDevice.Close()

	if err := hidDevice.SetAutoDetach(true); err != nil {
		logger.Error("unable to set auto detach", "err", err)
		return
	}
	if err := hidDevice.SetTarget(1, 0, 0); err != nil {
		logger.Error("unable to set target of hid device", "err", err)
		return
	}

	// Get report descriptor
	desc, err := hidDevice.GetReportDescriptor()
	if err != nil {
		logger.Error("unable to get report descriptor", "err", err)
		return
	}
	str, err := desc.String()
	if err != nil {
		logger.Error("unable to get descriptor string", "err", err)
		return
	}
	fmt.Printf("%s\n", str)
	fmt.Printf("%v\n", desc)

	// Create new Ethereum app instance
	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)
	ethApp := eth.NewEthereumApp(adpuProto, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	walletPath := "m'/44'/60'/0'/0/0"

	eip712Sig, err := ethApp.SignEIP712Message(ctx, walletPath, eip712.Message{
		Types: eip712.TypeStructs{
			{
				Name: "Transfer",
				Members: []eip712.FieldDefinition{
					{
						// from address;
						TypeDescription: eip712.FieldTypeDescription{
							IsArray:         false,
							IsSizeSpecified: false,
							Type:            eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "from",
					},
					{
						// to address;
						TypeDescription: eip712.FieldTypeDescription{
							IsArray:         false,
							IsSizeSpecified: false,
							Type:            eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "to",
					},
					{
						// ttt address;
						TypeDescription: eip712.FieldTypeDescription{
							IsArray:         false,
							IsSizeSpecified: false,
							Type:            eip712.FIELD_TYPE_DESC_TYPE_ADDRESS,
						},
						KeyName: "ttt",
					},
					{
						// amount uint64;
						TypeDescription: eip712.FieldTypeDescription{
							IsArray:         false,
							IsSizeSpecified: true,
							Type:            eip712.FIELD_TYPE_DESC_TYPE_UINT,
						},
						TypeSize: 8,
						KeyName:  "amount",
					},
				},
			},
		},
		Domain: eip712.Domain{
			Name:    "ERC20Transfer",
			Version: "0.1.0",
			ChainID: uint256.NewInt(10),
			VerifyingContract: schema.Address{
				0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45,
			},
		},
		Primary: eip712.StructItem{
			TypeName: "Transfer",
			Members: []eip712.StructItemMember{
				{
					Name: "from",
					Item: eip712.AtomicItem{
						Item: eip712.AddressData{
							0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
							0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
						},
					},
				},
				{
					Name: "to",
					Item: eip712.AtomicItem{
						Item: eip712.AddressData{
							0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
							0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
						},
					},
				},
				{
					Name: "ttt",
					Item: eip712.AtomicItem{
						Item: eip712.AddressData{
							0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
							0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42,
						},
					},
				},
				{
					Name: "amount",
					Item: eip712.AtomicItem{
						Item: eip712.NumberData{
							Num:     uint256.NewInt(123456),
							NumBits: 64,
							Signed:  false,
						},
					},
				},
			},
		},
		ClearSigning: eip712.ClearSigning{
			Enabled: false,
		},
	})

	if err != nil {
		logger.Error("unable to sign EIP712 message", "err", err)
		return
	}

	logger.Info("Signature", "R", log.HexDisplay(eip712Sig.R[:]), "S", log.HexDisplay(eip712Sig.S[:]), "V", eip712Sig.V)
}
