package internal

import (
	"context"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func FindDevices(timeout time.Duration) ([]string, error) {
	data := make(map[string]bool)
	var results []string

	// set a timeout to call stopscan after duration
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // cancel defers so if we stop early for some reason the goroutine below stops
	go func() {
		<-ctx.Done()
		globalAdapter.StopScan()
	}()

	globalAdapter.Enable()
	err := globalAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if strings.Contains(strings.ToLower(result.LocalName()), "aranet4") {
			_, exists := data[result.Address.String()]
			if !exists {
				data[result.Address.String()] = true // dedupe on address
				results = append(results, result.Address.String())
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
