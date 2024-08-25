package schema

import "encoding/binary"

type ETH2PublicKey struct {
	RawResponse
}

type ETH2WithdrawalIndex uint32

func (r *ETH2WithdrawalIndex) MarshalADPU() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(*r))

	return buf, nil
}
