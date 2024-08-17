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
)

func main() {
	logger := slog.Default()
	usbCtx := usb.NewGOUSBContext()
	man := manager.NewDeviceManager(usbCtx)

	defer func() {
		if err := man.Close(); err != nil {
			panic(err)
		}
	}()

	// List USB HID devices
	deviceInfos, err := man.Enumerate(0, 0)
	if err != nil {
		panic(err)
	}

	logger.Info("Device info", "info", deviceInfos.String())

	// Open Ledger Nano S
	hidDevice, err := man.Open(0x2C97, 0x1015, hid.DeviceConfig{
		StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	})
	if err != nil {
		panic(err)
	}
	defer hidDevice.Close()

	if err := hidDevice.SetTarget(1, 0, 0); err != nil {
		panic(err)
	}

	// Get report descriptor
	desc, err := hidDevice.GetReportDescriptor()
	if err != nil {
		panic(err)
	}

	str, err := desc.String()
	if err != nil {
		panic(err)
	}
	logger.Info("Report descriptor", "desc", str)

	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)

	ctx := context.Background()
	res, sw, err := adpuProto.Send(ctx, 0xE0, 0x06, 0x00, 0x00, nil)
	if err != nil {
		panic(err)
	}

	logger.Info("SW", "sw", sw)
	if sw != 0x9000 {
		panic(fmt.Errorf("sw is not OK"))
	}

	logger.Info("Ledger ETH information", "arbitraryDataEnabled", res[0]&0x01, "erc20ProvisioningNecessary", res[0]&0x02, "starkEnabled", res[0]&0x04, "starkv2Supported", res[0]&0x08, "version", fmt.Sprintf("%d.%d.%d", res[1], res[2], res[3]))
}
