package readings

type DeviceStatus uint8

const (
	A4StatusGreen DeviceStatus = 1
	A4StatusAmber DeviceStatus = 2
	A4StatusRed   DeviceStatus = 3
)

type DeviceReadings struct {
	CO2         uint16
	Temperature uint16
	Pressure    uint16
	Humidity    uint8
	Battery     uint8
	Status      DeviceStatus
}
