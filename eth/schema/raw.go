package schema

type RawResponse []byte

func (r *RawResponse) UnmarshalADPU(data []byte) error {
	*r = make([]byte, len(data))
	copy(*r, data)

	return nil
}

type EmptyResponse struct{}

func (r *EmptyResponse) UnmarshalADPU(data []byte) error {
	return nil
}
