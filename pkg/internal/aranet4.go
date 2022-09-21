package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KeganHollern/go-aranet4/pkg/aranet4/readings"
	"tinygo.org/x/bluetooth"
)

var (
	globalAdapter = bluetooth.DefaultAdapter
)

const (
	// Aranet UUIDs and handles
	// Services
	AR4_SERVICE     = "f0cd1400-95da-4f4b-9ac8-aa55d312af0c"
	GENERIC_SERVICE = "00001800-0000-1000-8000-00805f9b34fb"
	COMMON_SERVICE  = "0000180a-0000-1000-8000-00805f9b34fb"

	// Read / Aranet service
	AR4_READ_CURRENT_READINGS     = "f0cd1503-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_CURRENT_READINGS_DET = "f0cd3001-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_INTERVAL             = "f0cd2002-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_SECONDS_SINCE_UPDATE = "f0cd2004-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_TOTAL_READINGS       = "f0cd2001-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_HISTORY_READINGS_V1  = "f0cd2003-95da-4f4b-9ac8-aa55d312af0c"
	AR4_READ_HISTORY_READINGS_V2  = "f0cd2005-95da-4f4b-9ac8-aa55d312af0c"

	// Read / Generic servce
	GENERIC_READ_DEVICE_NAME = "00002a00-0000-1000-8000-00805f9b34fb"

	// Read / Common servce
	COMMON_READ_MANUFACTURER_NAME = "00002a29-0000-1000-8000-00805f9b34fb"
	COMMON_READ_MODEL_NUMBER      = "00002a24-0000-1000-8000-00805f9b34fb"
	COMMON_READ_SERIAL_NO         = "00002a25-0000-1000-8000-00805f9b34fb"
	COMMON_READ_HW_REV            = "00002a27-0000-1000-8000-00805f9b34fb"
	COMMON_READ_SW_REV            = "00002a28-0000-1000-8000-00805f9b34fb"
	COMMON_READ_BATTERY           = "00002a19-0000-1000-8000-00805f9b34fb"

	// Write / Aranet service
	AR4_WRITE_CMD = "f0cd1402-95da-4f4b-9ac8-aa55d312af0c"
)

type Aranet4Device struct {
	Mac string

	scan     bluetooth.Addresser
	device   *bluetooth.Device
	services []bluetooth.DeviceService
}

func (dev *Aranet4Device) ScanForDevice() error {
	globalAdapter.Enable()

	// set a timeout to call stopscan after duration
	scanning := true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel() // cancel defers so if we stop early for some reason the goroutine below stops
	go func() {
		<-ctx.Done()
		if scanning {
			globalAdapter.StopScan()
		}
	}()

	err := globalAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == dev.Mac {
			dev.scan = result.Address
			err := adapter.StopScan()
			if err != nil {
				panic(err)
			}
			scanning = false
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

func (dev *Aranet4Device) Connect() error {
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
func (dev *Aranet4Device) Address() string {
	if dev.Mac == "" {
		return dev.scan.String()
	}
	return dev.Mac
}

func (dev *Aranet4Device) DumpDevice() error {
	fmt.Println("----- DUMP ------")
	svcs, err := dev.device.DiscoverServices(nil)
	buf := make([]byte, 255)
	if err != nil {
		return err
	}
	for _, svc := range svcs {
		fmt.Printf("- service: %s\n", svc.UUID().String())
		chars, err := svc.DiscoverCharacteristics(nil)
		if err != nil {
			fmt.Printf("ERR: %s\n", err.Error())
		}
		for _, char := range chars {
			fmt.Printf("-- characterstic: %s\n", char.UUID().String())
			n, err := char.Read(buf)
			if err != nil {
				fmt.Printf("    ERR: %s\n", err.Error())
			} else {
				fmt.Printf("    len: %d\n", n)
				fmt.Printf("    val: %s\n", string(buf[:n]))
			}
		}
	}
	fmt.Println("----- DONE ------")
	return nil
}

// ------------------------ ar4 read operations & parsing

func (dev *Aranet4Device) Current(detailed bool) (*readings.DeviceReadings, error) {

	svc, err := dev.getService(AR4_SERVICE)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, err
	}
	char, err := dev.getCharacteristic(svc, AR4_READ_CURRENT_READINGS) // use AR4_READ_CURRENT_READINGS_DET to include interval & ago
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 255)
	n, err := char.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != 9 && n != 13 {
		// ????
		return nil, errors.New("malformatted response from device")
	}
	data := &readings.DeviceReadings{}
	data.Decode(buf, n) // supports CURRENT_READINGS and CURRENT_READINGS_DET

	return data, nil
}

// ------------------------ helpers for ar4 UUIDs (todo format uuids into byte[16])

func (dev *Aranet4Device) getCharacteristic(svc *bluetooth.DeviceService, a4_uuid string) (*bluetooth.DeviceCharacteristic, error) {
	if dev.device == nil {
		return nil, errors.New("not connected")
	}
	if svc == nil {
		return nil, errors.New("invalid service")
	}
	chars, err := svc.DiscoverCharacteristics(nil)
	if err != nil {
		return nil, err
	}
	for _, char := range chars {
		if char.UUID().String() == a4_uuid {
			return &char, nil
		}
	}
	return nil, errors.New("characteristic not found")
}
func (dev *Aranet4Device) getService(a4_uuid string) (*bluetooth.DeviceService, error) {
	if dev.device == nil {
		return nil, errors.New("not connected")
	}

	// if no services previously queried, lets get them all
	if len(dev.services) == 0 {
		var err error
		dev.services, err = dev.device.DiscoverServices(nil)
		if err != nil {
			return nil, err
		}
	}

	for _, svc := range dev.services {
		if svc.UUID().String() == a4_uuid {
			return &svc, nil
		}
	}

	return nil, errors.New("service not found")
}
