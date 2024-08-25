package schema

import (
	"encoding/binary"
	"fmt"

	"github.com/ntchjb/ledger-go/adpu"
	"github.com/ntchjb/ledger-go/eth/rlp"
)

type TxType uint8

const (
	// Legacy transcation type
	//
	// unsigned: rlp[nonce, gasprice, startgas, to, value, data, chainid, 0, 0]
	//
	// signing content: keccak256(rlp[nonce, gasprice, startgas, to, value, data, chainid, 0, 0])
	//
	// signed tx: rlp[nonce, gasprice, startgas, to, value, data, v, r, s]
	TX_TYPE_LEGACY TxType = 0x00
	// EIP-2930 transaction type
	//
	// unsigned: 0x01 || rlp([chainId, nonce, gasPrice, gasLimit, to, value, data, accessList])
	//
	// signing content: keccak256(unsigned)
	//
	// signed tx: 0x01 || rlp([chainId, nonce, gasPrice, gasLimit, to, value, data, accessList, signatureYParity, signatureR, signatureS])
	TX_TYPE_ACCESS_LIST TxType = 0x01
	// EIP-1559 transaction type
	//
	// unsigned: 0x02 || rlp([chain_id, nonce, max_priority_fee_per_gas, max_fee_per_gas, gas_limit, destination, amount, data, access_list])
	//
	// signing content: keccak256(unsigned)
	//
	// signed tx: 0x02 || rlp([chain_id, nonce, max_priority_fee_per_gas, max_fee_per_gas, gas_limit, destination, amount, data, access_list, signature_y_parity, signature_r, signature_s])
	TX_TYPE_DYNAMIC_FEE TxType = 0x02
	// EIP-4844 transaction type
	//
	// unsigned: 0x03 || rlp([chain_id, nonce, max_priority_fee_per_gas, max_fee_per_gas, gas_limit, to, value, data, access_list, max_fee_per_blob_gas, blob_versioned_hashes])
	//
	// signing content: keccak256(unsigned)
	//
	// signed tx: 0x03|| rlp([chain_id, nonce, max_priority_fee_per_gas, max_fee_per_gas, gas_limit, to, value, data, access_list, max_fee_per_blob_gas, blob_versioned_hashes, y_parity, r, s])
	TX_TYPE_BLOB TxType = 0x03
)

var (
	EIP2718TransactionTypes []bool = []bool{
		TX_TYPE_ACCESS_LIST: true,
		TX_TYPE_DYNAMIC_FEE: true,
		TX_TYPE_BLOB:        true,
	}

	SupportedTxTypes []bool = []bool{
		TX_TYPE_LEGACY:      true,
		TX_TYPE_ACCESS_LIST: true,
		TX_TYPE_DYNAMIC_FEE: true,
		TX_TYPE_BLOB:        false,
	}
)

type ChainID uint64

type TxInfo struct {
	// Transaction type
	TxType TxType
	// Transaction data payload i.e. calldata
	Data []byte
	// Target address used by this tx
	To Address
	// Chain ID i.e. ethereum = 0x01
	ChainID ChainID
	// Beginning position of chain ID data
	// This will be used to mitigate Ledger bug
	ChainIDOffset int
}

func DecodeTxInfo(rawTx []byte) (TxInfo, error) {
	var txInfo TxInfo
	// For Legacy Tx, this byte is >=0xC0 due to RLP encoding
	txType := TX_TYPE_LEGACY
	if int(rawTx[0]) < len(EIP2718TransactionTypes) && EIP2718TransactionTypes[rawTx[0]] {
		txType = TxType(rawTx[0])
	}
	rlpPart := rawTx
	if txType != TX_TYPE_LEGACY {
		rlpPart = rlpPart[1:]
	}

	rlpItem, n, err := rlp.Decode(rlpPart)
	if err != nil {
		return txInfo, fmt.Errorf("unable to decode RLP data from raw tx: %w", err)
	}
	if n < len(rlpPart) {
		return txInfo, fmt.Errorf("incomplete RLP data decoding, expected %d, but got %d decoded", len(rlpPart), n)
	}

	var data []byte
	var to Address
	var chainID ChainID
	switch txType {
	case TX_TYPE_BLOB:
		fallthrough
	case TX_TYPE_DYNAMIC_FEE:
		data = rlpItem.List[7].Data
		to = Address(rlpItem.List[5].Data)
		chainID = ChainID(rlpItem.List[0].Uint64())
	case TX_TYPE_ACCESS_LIST:
		data = rlpItem.List[6].Data
		to = Address(rlpItem.List[4].Data)
		chainID = ChainID(rlpItem.List[0].Uint64())
	default:
		data = rlpItem.List[5].Data
		to = Address(rlpItem.List[3].Data)
		if len(rlpItem.List) > 6 {
			chainID = ChainID(rlpItem.List[6].Uint64())
		} else {
			// For non EIP-155 transaction
			chainID = 1
		}
	}

	chainIDOffset := 0
	if txType == TX_TYPE_LEGACY && len(rlpItem.List) > 6 {
		if len(rlpItem.List) > 6 {
			txData := rlpItem
			txData.List = txData.List[len(txData.List)-3:]
			last3ItemsLength := txData.List[0].Len() + txData.List[1].Len() + txData.List[2].Len()
			chainIDOffset = len(rawTx) - last3ItemsLength

		} else {
			chainIDOffset = len(rawTx)
		}
	}

	txInfo.TxType = txType
	txInfo.ChainID = chainID
	txInfo.To = to
	txInfo.Data = data
	txInfo.ChainIDOffset = chainIDOffset

	return txInfo, nil
}

type SignTxRequest struct {
	// HD wallet path used for signing
	BIP32Path BIP32Path
	// RLP serialized transaction data to be signed
	Data []byte
}

func (c *SignTxRequest) MarshalADPU() ([]byte, error) {
	buf, err := adpu.Marshal(&c.BIP32Path)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal BIP-32 path: %w", err)
	}

	buf = append(buf, c.Data...)

	return buf, nil
}

type SignPersonalMessageRequest struct {
	// HD wallet path used for signing
	BIP32Path BIP32Path
	// Personal message, can have maximum length of MAX_UINT32
	Data []byte
}

func (c *SignPersonalMessageRequest) MarshalADPU() ([]byte, error) {
	buf, err := adpu.Marshal(&c.BIP32Path)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal BIP-32 path: %w", err)
	}

	messageLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(messageLengthBytes, uint32(len(c.Data)))
	buf = append(buf, messageLengthBytes...)
	buf = append(buf, c.Data...)

	return buf, nil
}

type SignEIP712HashedRequest struct {
	// HD wallet path used for signing
	BIP32Path             BIP32Path
	HashedDomainSeparator [32]byte
	HashedMessage         [32]byte
}

func (r *SignEIP712HashedRequest) MarshalADPU() ([]byte, error) {
	buf, err := adpu.Marshal(&r.BIP32Path)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal BIP-32 path: %w", err)
	}
	buf = append(buf, r.HashedDomainSeparator[:]...)
	buf = append(buf, r.HashedMessage[:]...)

	return buf, nil
}

type SignatureV uint64

// Recover V value for chain with large chain ID
// Because Ledger can only send 1 byte of V value to client,
// but EIP-155 tx requires V = y_parity + chain_id * 2 + 35
// which can be larger than 1 byte (>255).
// To mitigate this, we can extract Y parity, can re-calculate V value.
func (v SignatureV) RecoverLegacy(chainID ChainID) SignatureV {
	// If result of highest V value possible cannot fit in 1 byte,
	// then get Y parity from V value returned by Ledger first,
	// and try re-calculate V value from the Y parity and chain ID.
	if chainID*2+35+1 > 255 {
		// Simulate what Ledger device calculates
		// by getting highest 4 bytes of chainID
		chainIDTruncated := chainID >> 32
		// V value returned from Ledger can be either
		// - (chain_id * 2 + 35 + 0) % 256
		// - (chain_id * 2 + 35 + 1) % 256
		ledgerVValueLowest := uint8((chainIDTruncated*2 + 35) % 256)
		// This is actual Y parity
		yParity := Abs(uint8(v) - ledgerVValueLowest)

		v = SignatureV(uint64(chainID)*2 + 35 + uint64(yParity))
	}

	return v
}

type SignDataResponse struct {
	V SignatureV
	R [32]byte
	S [32]byte
}

func (r *SignDataResponse) UnmarshalADPU(data []byte) error {
	if len(data) < 65 {
		return fmt.Errorf("data is too short, expected 65, got %d", len(data))
	}

	r.V = SignatureV(data[0])
	copy(r.R[:], data[1:33])
	copy(r.S[:], data[33:65])

	return nil
}
