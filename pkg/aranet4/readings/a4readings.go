package readings

import (
	"encoding/binary"
	"errors"
)

type DeviceStatus uint8

const (
	A4StatusGreen DeviceStatus = 1
	A4StatusAmber DeviceStatus = 2
	A4StatusRed   DeviceStatus = 3
)

type DeviceReadingsExt struct {
	Interval uint16
	Ago      uint16
}
type DeviceReadings struct {
	CO2         uint16
	Temperature uint16
	Pressure    uint16
	Humidity    uint8
	Battery     uint8
	Status      DeviceStatus

	DeviceReadingsExt
}

func (data *DeviceReadings) Decode(buf []byte, len int) error {
	if len < 9 || len > 13 {
		return errors.New("invalid length unable to decode")
	}
	if len >= 9 {
		data.CO2 = binary.LittleEndian.Uint16(buf[0:2])
		data.Temperature = binary.LittleEndian.Uint16(buf[2:4])
		data.Pressure = binary.LittleEndian.Uint16(buf[4:6])
		data.Humidity = buf[6]
		data.Battery = buf[7]
		data.Status = DeviceStatus(buf[8])
		data.Interval = 0
		data.Ago = 0
	}
	if len == 13 {
		data.Interval = binary.LittleEndian.Uint16(buf[9:11])
		data.Ago = binary.LittleEndian.Uint16(buf[11:13])
	}
	return nil
}
