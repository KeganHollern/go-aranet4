package internal

import (
	"errors"
	"strings"

	"tinygo.org/x/bluetooth"
)

var (
	globalAdapter = bluetooth.DefaultAdapter
)

type Aranet4Device struct {
	Mac string

	scan   bluetooth.Addresser
	device *bluetooth.Device
}

func FindDevices() ([]string, error) {
	var results []string
	globalAdapter.Enable()
	err := globalAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if strings.Contains(strings.ToLower(result.LocalName()), "aranet4") {
			results = append(results, result.Address.String())
		}
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (dev Aranet4Device) ScanForDevice() error {
	globalAdapter.Enable()

	err := globalAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == dev.Mac {
			dev.scan = result.Address
			err := adapter.StopScan()
			if err != nil {
				panic(err)
			}
		}
	})
	if err != nil {
		return err
	}
	if dev.scan == nil {
		return errors.New("device not found")
	}
	return nil
}

func (dev Aranet4Device) Connect() error {
	if dev.scan == nil {
		err := dev.ScanForDevice() // scan again ?
		if err != nil {
			return err
		}
	}

	var err error
	dev.device, err = globalAdapter.Connect(dev.scan, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	if dev.device == nil {
		return errors.New("failed to connect")
	}
	return nil
}
func (dev Aranet4Device) Disconnect() {
	if dev.device != nil {
		dev.device.Disconnect()
	}
}
