package aranet4

import (
	"errors"
	"github.com/KeganHollern/go-aranet4/pkg/internal"
)

// A4Device defines an interface for interacting with Aranet4 Bluetooth Devices.
type A4Device interface {
	Connect() error
}

// NewDeviceFromMAC returns an A4Device from the provided bluetooth MAC address.
// If the provided MAC is not found (or paired), an error is returned.
func NewDeviceFromMAC(MAC string) (A4Device, error) {

	device := &internal.Aranet4Device{}

	return nil, errors.New("not implemented")
}

// GetNearbyDevices returns a list of A4Device for each nearby Aranet4 bluetooth device.
// Note that only paired Aranet4 devices are returned.
func GetNearbyDevices() ([]A4Device, error) {
	return nil, errors.New("not implemented")
}
