package schema

type RawRequest []byte

func (r *RawRequest) MarshalADPU() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	return *r, nil
}

type RawResponse []byte

func (r *RawResponse) UnmarshalADPU(data []byte) error {
	*r = make([]byte, len(data))
	copy(*r, data)

	return nil
}

type EmptyRequest struct{}

func (r *EmptyRequest) MarshalADPU() ([]byte, error) {
	return []byte{}, nil
}

type EmptyResponse struct{}

func (r *EmptyResponse) UnmarshalADPU(data []byte) error {
	return nil
}
