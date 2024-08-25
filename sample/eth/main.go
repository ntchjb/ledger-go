package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"time"

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

	// Create new Ethereum app instance
	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)
	ethApp := eth.NewEthereumApp(adpuProto, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	walletPath := "m'/44'/60'/0'/0/0"

	// #1 ETH: get configuration
	conf, err := ethApp.GetConfiguration(ctx)
	if err != nil {
		logger.Error("unable to get configuration", "err", err)
		return
	}
	logger.Info("Configuration", "conf", conf)

	// #2 ETH: get address
	address, err := ethApp.GetAddress(ctx, walletPath, false, false, 1)
	if err != nil {
		logger.Error("unable to get address", "err", err)
		return
	}
	logger.Info("Address", "addr", address.Address.String(), "publicKey", log.HexDisplay(address.PublicKey[:]), "chaincode", log.HexDisplay(address.Chaincode[:]))

	// #3 ETH: sign transaction
	// Legacy
	// rawTx, _ := hex.DecodeString("f901e980830f4240830205b5943cd1dfb81c50a5300c60a181ed145a7286d81e0a80b901c4183fb413000000000000000000000000794a61358d6845594f94dc1db02a252b5b4814ad0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000000000000000000b41544f4b454e5f494d504c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b41544f4b454e5f494d504c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000430783030000000000000000000000000000000000000000000000000000000000a0000")
	// EIP-1559
	rawTx, _ := hex.DecodeString("02f870018313fc97808432c3453a825a3c94388c818ca8b9251b393131c08a736a67ccb192978768f233feb2c98a80c080a0e21a0b9a80dc27cd2c9ccc551a7df692b83d2a522aa62fd47949f07363afcceaa07aaef211074d6e8c132e937202da0a0ce6648f328cd6d5e90e41b82955e3b224")
	txSig, err := ethApp.SignTransaction(ctx, walletPath, rawTx)
	if err != nil {
		logger.Error("unable to sign tx", "err", err)
		return
	}
	logger.Info("Signature", "R", log.HexDisplay(txSig.R[:]), "S", log.HexDisplay(txSig.S[:]), "V", txSig.V)

	// #4 ETH: sign personal message
	pmSig, err := ethApp.SignPersonalMessage(ctx, walletPath, []byte("hello world"))
	if err != nil {
		logger.Error("unable to sign personal message", "err", err)
		return
	}
	logger.Info("Signature", "R", log.HexDisplay(pmSig.R[:]), "S", log.HexDisplay(pmSig.S[:]), "V", pmSig.V)

	// #5 ETH: sign EIP-712 hashed message
	domainSeparator, _ := hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
	message, _ := hex.DecodeString("0202020202020202020202020202020202020202020202020202020202020202")
	eip712Sig, err := ethApp.SignEIP712MessageHash(ctx, walletPath, domainSeparator, message)
	if err != nil {
		logger.Error("unable to sign EIP712 hashed message", "err", err)
		return
	}
	logger.Info("Signature", "R", log.HexDisplay(eip712Sig.R[:]), "S", log.HexDisplay(eip712Sig.S[:]), "V", eip712Sig.V)

	// #6 ETH: ETH2 get public key
	publicKey, err := ethApp.ETH2GetPublicKey(ctx, walletPath, true)
	if err != nil {
		logger.Error("unable to get public key for ETH2", "err", err)
		return
	}
	logger.Info("ETH2 public key", "publicKey", log.HexDisplay(publicKey.RawResponse))

	// #7 ETH: ETH2 set withdrawal index
	if err := ethApp.ETH2SetWithdrawalIndex(ctx, 0); err != nil {
		logger.Error("unable to set withdrawal index for ETH2", "err", err)
		return
	}

	// #8: ETH: Get privacy public key and shared secret (X25519)
	privacyPublicKey, err := ethApp.GetPrivacyPublicKey(ctx, walletPath, false)
	if err != nil {
		logger.Error("unable to get privacy public key", "err", err)
		return
	}
	logger.Info("ed25519 public key", "publicKey", log.HexDisplay(privacyPublicKey.RawResponse))

	remotePublicKey, _ := hex.DecodeString("87020e80af6e07a6e4697f091eacadb9e7e6629cb7e5a8a371689a3ed53b3d64")
	privacySharedSecret, err := ethApp.GetPrivacySharedSecret(ctx, walletPath, remotePublicKey, false)
	if err != nil {
		logger.Error("unable to get privacy public key", "err", err)
		return
	}
	logger.Info("ed25519 shared secret", "sharedSecret", log.HexDisplay(privacySharedSecret.RawResponse))
}
