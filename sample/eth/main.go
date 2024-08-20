package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ntchjb/gohid/hid"
	"github.com/ntchjb/gohid/manager"
	"github.com/ntchjb/gohid/usb"
	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/device"
	"github.com/ntchjb/ledger-go/eth"
)

func main() {
	logger := slog.Default()
	usbCtx := usb.NewGOUSBContext()
	man := manager.NewDeviceManager(usbCtx, logger)

	defer func() {
		if err := man.Close(); err != nil {
			logger.Error("unable to close device manager", "err", err)
		}
	}()

	// List USB HID devices
	deviceInfos, err := man.Enumerate(0, 0)
	if err != nil {
		logger.Error("unable to enumerate devices", "err", err)
		return
	}

	fmt.Printf("DeviceInfos:\n%s\n", deviceInfos)

	// Open Ledger Nano S
	hidDevice, err := man.Open(0x2C97, 0x1011, hid.DeviceConfig{
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

	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)
	ethApp := eth.NewEthereumApp(adpuProto)

	ctx := context.Background()

	conf, err := ethApp.GetConfiguration(ctx)
	if err != nil {
		logger.Error("unable to get configuration", "err", err)
		return
	}
	logger.Info("Configuration", "conf", conf)

	address, err := ethApp.GetAddress(ctx, "m'/44'/60'/0'/0/0", false, false, 1)
	if err != nil {
		logger.Error("unable to get address", "err", err)
		return
	}
	logger.Info("Address", "addr", address.Address.String(), "publicKey", address.PublicKey.String(), "chaincode", address.Chaincode.String())
}
