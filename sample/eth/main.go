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
			logger.Error("unable to close device manager", "err", err)
		}
	}()

	// List USB HID devices
	deviceInfos, err := man.Enumerate(0, 0)
	if err != nil {
		logger.Error("unable to enumerate devices", "err", err)
		return
	}

	logger.Info("Device info", "info", deviceInfos.String())

	// Open Ledger Nano S
	hidDevice, err := man.Open(0x2C97, 0x1015, hid.DeviceConfig{
		StreamLaneCount: hid.DEFAULT_ENDPOINT_STREAM_COUNT,
	})
	if err != nil {
		logger.Error("unable to open device")
		return
	}
	defer hidDevice.Close()

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

	ctx := context.Background()
	res, sw, err := adpuProto.Send(ctx, 0xE0, 0x06, 0x00, 0x00, nil)
	if err != nil {
		logger.Error("unable to send ADPU protocol to device", "err", err)
	}

	logger.Info("SW", "sw", sw)
	if sw != 0x9000 {
		logger.Error("SW is not OK")
		return
	}

	logger.Info("Ledger ETH information", "arbitraryDataEnabled", res[0]&0x01, "erc20ProvisioningNecessary", res[0]&0x02, "starkEnabled", res[0]&0x04, "starkv2Supported", res[0]&0x08, "version", fmt.Sprintf("%d.%d.%d", res[1], res[2], res[3]))
}
