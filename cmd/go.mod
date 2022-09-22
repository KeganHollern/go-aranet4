module github.com/KeganHollern/go-aranet4/cmd

go 1.19

require (
	github.com/KeganHollern/go-aranet4/pkg v0.0.1
	github.com/sirupsen/logrus v1.9.0
)

replace github.com/KeganHollern/go-aranet4/pkg v0.0.1 => ../pkg

require (
	github.com/JuulLabs-OSS/cbgo v0.0.2 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/godbus/dbus/v5 v5.0.3 // indirect
	github.com/muka/go-bluetooth v0.0.0-20210812063148-b6c83362e27d // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	tinygo.org/x/bluetooth v0.5.0 // indirect
)
