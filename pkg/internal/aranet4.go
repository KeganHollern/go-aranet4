package internal

import "tinygo.org/x/bluetooth"

var (
	globalAdapter = bluetooth.DefaultAdapter
)

type Aranet4Device struct {
	Mac string
}

func (dev Aranet4Device) Connect() error {
	globalAdapter.Enable()
	var foundDevice bluetooth.ScanResult
	err := globalAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == dev.Mac {
			foundDevice = result
			err := adapter.StopScan()
			if err != nil {
				panic(err) // unrecoverable
			}
		}
	})
	if err != nil {
		return err
	}

	device, err := globalAdapter.Connect()
}
