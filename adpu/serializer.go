package adpu

type Marshaler interface {
	MarshalADPU() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalADPU(data []byte) error
}

func Marshal[T Marshaler](m T) ([]byte, error) {
	return m.MarshalADPU()
}

func Unmarshal[T Unmarshaler](data []byte, target T) error {
	return target.UnmarshalADPU(data)
}

type EmptyData struct{}

func (e *EmptyData) MarshalADPU() ([]byte, error) {
	return nil, nil
}

func (e *EmptyData) UnmarshalADPU(data []byte) error {
	return nil
}
