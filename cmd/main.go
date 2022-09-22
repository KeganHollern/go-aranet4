package main

import (
	"time"

	"github.com/KeganHollern/go-aranet4/pkg/aranet4"
	"github.com/KeganHollern/go-aranet4/pkg/aranet4/types"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Infoln("starting go-aranet4 sample")

	devices, err := aranet4.GetNearbyDevices(time.Second * 5)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to get nearby devices")
	}

	// list all found aranet4 devices
	for _, device := range devices {
		addr := device.Address()
		logrus.WithField("mac", addr).Infoln("found aranet4 device")
		continue

	}
	if len(devices) == 0 {
		logrus.Warnln("no devices found")
		return // exit early we found no devices
	}
	// we'll use the first device for the demo
	device := devices[0]

	// connect to the device
	logrus.WithField("mac", device.Address()).Infoln("connecting to device...")
	err = device.Connect()
	if err != nil {
		logrus.WithError(err).Errorln("failed to connect to device")
		return
	}
	defer device.Disconnect() // defer disconnection
	logrus.WithField("mac", device.Address()).Infoln("connected")

	// dump details
	readings, err := device.Current(true)
	if err != nil {
		logrus.WithError(err).Errorln("failed to read from device")
		return
	}
	logrus.WithField("readings", readings).Infoln("device read successfully")

	// get some accurate measurements
	cancel, err := device.Subscribe(func(readings *types.A4Data) {
		logrus.WithField("readings", readings).Infoln("subscription event")
	})
	if err != nil {
		logrus.WithError(err).Errorln("failed to subscribe to device")
	}

	time.Sleep(time.Minute * 6)
	cancel()
	time.Sleep(time.Minute * 6)
}
