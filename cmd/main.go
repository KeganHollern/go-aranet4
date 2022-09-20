package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugln("Starting Aranet4 Go Sample")
}
