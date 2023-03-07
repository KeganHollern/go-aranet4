package types

import (
	"encoding/binary"
	"errors"
)

type A4Status int

const (
	A4StatusGreen A4Status = 1
	A4StatusAmber A4Status = 2
	A4StatusRed   A4Status = 3
)

type A4DataExt struct {
	Interval int
	Ago      int
}
type A4Data struct {
	Raw         []byte
	CO2         int
	Temperature float64
	Pressure    float64
	Humidity    int
	Battery     int
	Status      A4Status

	A4DataExt
}

func (data *A4Data) Decode(buf []byte) error {
	n := len(buf)
	if n != 9 && n != 13 {
		return errors.New("invalid length unable to decode")
	}

	co2sensor := binary.LittleEndian.Uint16(buf[0:2])
	if (co2sensor & 0x8000) == 0x8000 {
		return errors.New("invalid co2 sensor value")
	}
	data.CO2 = int(co2sensor)

	tempsensor := binary.LittleEndian.Uint16(buf[2:4])
	if tempsensor == 0x4000 || tempsensor > 0x8000 {
		return errors.New("invalid temperature sensor value")
	}
	data.Temperature = float64(tempsensor) / 20.0
	pressensor := binary.LittleEndian.Uint16(buf[4:6])
	if (pressensor & 0x8000) == 0x8000 {
		return errors.New("invalid pressure sensor value")
	}
	data.Pressure = float64(pressensor) / 10.0
	data.Humidity = int(buf[6])
	if (data.Humidity & 0x80) == 0x80 {
		return errors.New("invalid humidity sensor value")
	}

	data.Battery = int(buf[7])
	data.Status = A4Status(buf[8])
	if n == 13 {
		data.Interval = int(binary.LittleEndian.Uint16(buf[9:11]))
		data.Ago = int(binary.LittleEndian.Uint16(buf[11:13]))
	}

	return nil
}
