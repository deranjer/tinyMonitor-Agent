module github.com/deranjer/tinyMonitor-Agent

go 1.12

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/deranjer/tinyMonitor v0.0.0-20190723015216-38f9e0c0d604
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/rs/zerolog v1.14.3
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/spf13/viper v1.4.0
	nanomsg.org/go/mangos/v2 v2.0.2
)

replace github.com/deranjer/tinyMonitor => ../tinyMonitor
