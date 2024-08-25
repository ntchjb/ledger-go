package schema

import (
	"fmt"

	"github.com/ntchjb/ledger-go/adpu"
)

type GetPrivacyPublicKeyResponse struct {
	RawResponse
}

type GetPrivacySharedSecretRequest struct {
	Path            BIP32Path
	RemotePublicKey []byte
}

func (r *GetPrivacySharedSecretRequest) MarshalADPU() ([]byte, error) {
	buf, err := adpu.Marshal(&r.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal BIP-32 path: %w", err)
	}

	buf = append(buf, r.RemotePublicKey...)

	return buf, nil
}

type GetPrivacySharedSecretResponse struct {
	RawResponse
}
