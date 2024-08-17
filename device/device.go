package device

import (
	"context"

	"github.com/ntchjb/gohid/hid"
)

type Device interface {
	Read(ctx context.Context, data []byte) (n int, err error)
	Write(ctx context.Context, data []byte) (n int, err error)
}

type ledgerDevice struct {
	device hid.Device
}

func NewLedgerDevice(hidDevice hid.Device) Device {
	return &ledgerDevice{
		device: hidDevice,
	}
}

func (d *ledgerDevice) Read(ctx context.Context, data []byte) (n int, err error) {
	return d.device.ReadInput(ctx, data)
}
func (d *ledgerDevice) Write(ctx context.Context, data []byte) (n int, err error) {
	// append report ID as 0, as Ledger doesn't use Report ID
	data = append([]byte{0x00}, data...)
	return d.device.WriteOutput(ctx, data)
}
