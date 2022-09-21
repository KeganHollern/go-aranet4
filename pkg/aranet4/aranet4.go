package aranet4

import (
	"github.com/KeganHollern/go-aranet4/pkg/internal"
)

// A4Device defines an interface for interacting with Aranet4 Bluetooth Devices.
type A4Device interface {
	Connect() error
	Disconnect()
}

// NewDeviceFromMAC returns an A4Device from the provided bluetooth MAC address.
// If the provided MAC is not found (or paired), an error is returned.
//
//	device, err := aranet4.NewDeviceFromMAC("00:00:00:00:00:00")
//	/* handle err */
//	err = device.Connect()
//	/* handle err */
//	defer device.Disconnect()
//	...
func NewDeviceFromMAC(MAC string) (A4Device, error) {
	a4 := &internal.Aranet4Device{
		Mac: MAC,
	}
	err := a4.ScanForDevice()
	if err != nil {
		return nil, err
	}
	return a4, nil
}

// GetNearbyDevices returns a list of A4Device for each nearby Aranet4 bluetooth device.
// Note that only paired Aranet4 devices are returned.
// Devices which fail to connect are dropped without error.
func GetNearbyDevices() ([]A4Device, error) {
	var results []A4Device
	macs, err := internal.FindDevices()
	if err != nil {
		return nil, err
	}

	for _, mac := range macs {
		device, err := NewDeviceFromMAC(mac)
		if err != nil {
			continue
		}
		results = append(results, device)
	}

	return results, nil
}
