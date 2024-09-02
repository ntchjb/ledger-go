package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
	"github.com/ntchjb/ledger-go/log"
	"github.com/ntchjb/ledger-go/sample/eth/util"
)

func GetDomainSignature(cli *http.Client, domainName string, challenge schema.Challenge) ([]byte, error) {
	url := fmt.Sprintf("https://nft.api.live.ledger.com/v1/names/ens/forward/%s?challenge=0x%x", domainName, challenge)
	resp, err := cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get domain signature: %w", err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	payload := schema.DomainSignatureResponse{}
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("unable to decode domain signature payload: %w", err)
	}

	payloadBytes, err := hex.DecodeString(payload.Payload)
	if err != nil {
		return nil, fmt.Errorf("unable to decode domain signature payload to bytes: %w", err)
	}

	return payloadBytes, nil
}

// This one requires 1inch app to be installed in Ledger Device
func get1inchOptimismResolutionAndPayload() (schema.ClearSigningResolution, []byte) {
	descSrcToken, _ := util.ERC20Sigs.FindByChainIDAndAddress(uint256.NewInt(10), schema.Address{
		0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06,
	})
	descDstToken, _ := util.ERC20Sigs.FindByChainIDAndAddress(uint256.NewInt(10), schema.Address{
		0x0b, 0x2C, 0x63, 0x9c, 0x53, 0x38, 0x13, 0xf4, 0xAa, 0x9D,
		0x78, 0x37, 0xCA, 0xf6, 0x26, 0x53, 0xd0, 0x97, 0xFf, 0x85,
	})
	oneInchSwapPayload, _ := hex.DecodeString("0531696e63681111111254eeb25477b68fb85ed929f73a96058212aa3caf")
	oneInchSwapSig, _ := hex.DecodeString("30450221009bf7192ed1276263000f619b6133c98a393bff309ac8901b5593849fbf276b2702202de029f07bd0573737b368a80d592daa331ad982b0b9162ecc301fb017846e8c")
	resolution := schema.ClearSigningResolution{
		ERC20Tokens: []schema.ERC20TokenResolution{
			// Derived from https://cdn.live.ledger.com/cryptoassets/evm/10/erc20-signatures.json
			// and referred to external plugin data below at `erc20OfInterest`
			descSrcToken.Raw,
			descDstToken.Raw,
		},
		ExternalPlugin: []schema.ExternalPluginResolution{
			{
				// Derived from https://cdn.live.ledger.com/plugins/ethereum.json
				// Contract: 0x1111111254eeb25477b68fb85ed929f73a960582
				// Selector: 0x12aa3caf
				Payload:   oneInchSwapPayload,
				Signature: oneInchSwapSig,
			},
		},
	}
	rawTx, _ := hex.DecodeString("f9034c83036988831ee50c830690eb941111111254eeb25477b68fb85ed929f73a96058280b9032412aa3caf000000000000000000000000b63aae6c353636d66df13b89ba4425cfe13d10ba00000000000000000000000042000000000000000000000000000000000000060000000000000000000000000b2c639c533813f4aa9d7837caf62653d097ff85000000000000000000000000b63aae6c353636d66df13b89ba4425cfe13d10ba0000000000000000000000003f343211f0487eb43af2e0e773ba012015e6651a0000000000000000000000000000000000000000000000000b59155ba7b59b8000000000000000000000000000000000000000000000000000000000773819f40000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018100000000000000000000000000000000000000000000000000000000016300a007e5c0d200000000000000000000000000000000000000000000000000013f00004f02a0000000000000000000000000000000000000000000000000000000000034727cee63c1e50185c31ffa3706d1cce9d525a00f1c7d4a2911754c420000000000000000000000000000000000000651204c4af8dbc524681930a27b2f1af5bcc8062e6fb768f180fcce6836688e9084f035309e29bf0a209500447dc2038200000000000000000000000068f180fcce6836688e9084f035309e29bf0a20950000000000000000000000000b2c639c533813f4aa9d7837caf62653d097ff8500000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000076154f3f0000000000000000000000001111111254eeb25477b68fb85ed929f73a96058200000000000000000000000042f527f50f16a103b6ccab48bccca214500c1021000000000000000000000000000000000000000000000000000000000000000a0000")

	return resolution, rawTx
}

func getBAYCResolutionAndPayload() (schema.ClearSigningResolution, []byte) {
	baycPluginData, _ := hex.DecodeString("010106455243373231bc4ca0eda7647a8ab7c2061c2e118a18a936f13d23b872dd0000000000000001020147304502204ab947bb134b9e42b0098a2292227e6042dc2c6dd3b7f65a9790b152c84dc1ac022100ccff370c40b8453b666a01220c8ea64558214acf876b8758e629901a944a407b")
	baycNFTData, _ := hex.DecodeString("010111426f7265644170655961636874436c7562bc4ca0eda7647a8ab7c2061c2e118a18a936f13d0000000000000001010147304502206987a6c26f8b42dfc4b410a5e192d218dadd2c7a76c367c8116a4806d21704460221008a2fab1950f9985ec6ef643e73f56b317f45e20996763c37008c39e32fbc51c5")
	resolution := schema.ClearSigningResolution{
		Plugin: []schema.PluginResolution{
			// Derived from https://nft.api.live.ledger.com/v1/ethereum/1/contracts/0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D/plugin-selector/0x23b872dd
			baycPluginData,
		},
		NFTs: []schema.NFTResolution{
			// Derived from https://nft.api.live.ledger.com/v1/ethereum/1/contracts/0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D
			baycNFTData,
		},
	}
	rawTx, _ := hex.DecodeString("02f88f018206f7841dcd650084682d1eae8302dec894bc4ca0eda7647a8ab7c2061c2e118a18a936f13d80b86423b872dd000000000000000000000000ec1c5f91ff6ca0351d0be13c88b5d9553ebc03a6000000000000000000000000fe89cc7abb2c4183683ab71653c4cdc9b02d44b7000000000000000000000000000000000000000000000000000000000000248bc0")

	return resolution, rawTx
}

func getTransferETHResolutionAndPayload() (schema.ClearSigningResolution, []byte) {
	resolution := schema.ClearSigningResolution{
		Domains: []schema.DomainResolution{
			{
				Registry: schema.DOMAIN_REGISTRY_ENS,
				Domain:   "vitalik.eth",
				Type:     schema.DOMAIN_TYPE_FORWARD,
			},
		},
	}
	rawTx, _ := hex.DecodeString("02e701188459682f008459682f0082520894d8da6bf26964af9d7eed9e03e53415d37aa960453480c0")

	return resolution, rawTx
}

func main() {
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	usbCtx := usb.NewGOUSBContext()
	man := manager.NewDeviceManager(usbCtx, logger)

	httpCli := &http.Client{
		Timeout: time.Second * 30,
	}

	defer func() {
		if err := man.Close(); err != nil {
			logger.Error("unable to close device manager", "err", err)
		}
	}()

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

	// Create new Ethereum app instance
	ledgerDevice := device.NewLedgerDevice(hidDevice)
	adpuProto := adpu.NewProtocol(ledgerDevice, 1234, logger)
	ethApp := eth.NewEthereumApp(adpuProto, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	walletPath := "m'/44'/60'/0'/0/0"

	// resolution, rawTx := get1inchOptimismResolutionAndPayload()
	// resolution, rawTx := getBAYCResolutionAndPayload()
	resolution, rawTx := getTransferETHResolutionAndPayload()
	for _, domain := range resolution.Domains {
		challengeData, err := ethApp.GetChallenge(ctx)
		if err != nil {
			logger.Error("unable to get challenge data", "err", err)
			return
		}
		logger.Info("ChallengeData", "challenge", log.HexDisplay(challengeData[:]))
		signature, err := GetDomainSignature(httpCli, domain.Domain, challengeData)
		if err != nil {
			logger.Error("unable to get domain signature", "err", err)
			return
		}

		logger.Info("domain signature", "sig", signature)

		if err := ethApp.ProvideDomainNameInformation(ctx, signature); err != nil {
			logger.Error("unable to provide domain name info", "err", err)
			return
		}
	}

	for _, plugin := range resolution.Plugin {
		if err := ethApp.SetPlugin(ctx, plugin); err != nil {
			logger.Error("unable to set plugin", "err", err)
			return
		}
	}

	for _, plugin := range resolution.ExternalPlugin {
		if err := ethApp.SetExternalPlugin(ctx, plugin.Payload, plugin.Signature); err != nil {
			logger.Error("unable to set external plugin", "err", err)
			return
		}
	}

	for _, nftSig := range resolution.NFTs {
		if err := ethApp.ProvideNFTInformation(ctx, nftSig); err != nil {
			logger.Error("unable to provide NFT information", "err", err)
			return
		}
	}

	for _, tokenSig := range resolution.ERC20Tokens {
		if _, err := ethApp.ProvideERC20Information(ctx, tokenSig); err != nil {
			logger.Error("unable to provide ERC20 token info", "err", err)
			return
		}
	}

	txSig, err := ethApp.SignTransaction(ctx, walletPath, rawTx)
	if err != nil {
		logger.Error("unable to sign tx", "err", err)
		return
	}
	logger.Info("Signature", "R", log.HexDisplay(txSig.R[:]), "S", log.HexDisplay(txSig.S[:]), "V", txSig.V)
}
