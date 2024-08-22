package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"github.com/ntchjb/gohid/hid"
	"github.com/ntchjb/gohid/manager"
	"github.com/ntchjb/gohid/usb"
	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/device"
	"github.com/ntchjb/ledger-go/eth"
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
	ethApp := eth.NewEthereumApp(adpuProto, logger)

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

	// Legacy
	// rawTx, _ := hex.DecodeString("f901e980830f4240830205b5943cd1dfb81c50a5300c60a181ed145a7286d81e0a80b901c4183fb413000000000000000000000000794a61358d6845594f94dc1db02a252b5b4814ad0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000000000000000000b41544f4b454e5f494d504c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b41544f4b454e5f494d504c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000430783030000000000000000000000000000000000000000000000000000000000a0000")
	// EIP-1559
	rawTx, _ := hex.DecodeString("02f870018313fc97808432c3453a825a3c94388c818ca8b9251b393131c08a736a67ccb192978768f233feb2c98a80c080a0e21a0b9a80dc27cd2c9ccc551a7df692b83d2a522aa62fd47949f07363afcceaa07aaef211074d6e8c132e937202da0a0ce6648f328cd6d5e90e41b82955e3b224")
	res, err := ethApp.SignTransaction(ctx, "m'/44'/60'/0'/0/0", rawTx)
	if err != nil {
		logger.Error("unable to sign tx", "err", err)
		return
	}

	logger.Info("Signature", "R", log.HexDisplay(res.R[:]), "S", log.HexDisplay(res.S[:]), "V", res.V)
}
