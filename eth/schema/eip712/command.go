package eip712

type CommandType uint8

const (
	COMMAND_TYPE_SEND_TYPE_DEFINITION CommandType = iota
	COMMAND_TYPE_SEND_DATA
	COMMAND_TYPE_SEND_FILTER
)

type DataCommand struct {
	Component Component
	Value     []byte
}
