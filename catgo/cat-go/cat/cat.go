package cat

import (
	"errors"
	"os"
	"sync/atomic"
)

var isEnabled uint32 = 0

type HostConfig struct {
	Host string
	Port int
}

func ResetLogger(l ILogger) {
	logger.logger = l
}

func InitWithServer(domain string, services []HostConfig) error {

	if err := config.Init(domain, false); err != nil {
		return err
	}

	for _, x := range services {
		config.httpServerAddresses = append(config.httpServerAddresses, serverAddress{
			host: x.Host,
			port: x.Port,
		})
	}

	if !enable() {
		return errors.New("cannot init twice")
	}

	logger.Info("Cat has been enabled.")

	go background(&router)
	go background(&monitor)
	go background(&sender)
	aggregator.Background()
	return nil

}

func Init(domain string) {
	if err := config.Init(domain, true); err != nil {
		logger.Warning("Cat initialize failed.")
		return
	}
	if !enable() {
		return
	}
	logger.Info("Cat has been enabled.")

	go background(&router)
	go background(&monitor)
	go background(&sender)
	aggregator.Background()
}

func enable() bool {
	return atomic.SwapUint32(&isEnabled, 1) == 0
}

func disable() {
	if atomic.SwapUint32(&isEnabled, 0) == 1 {
		logger.Info("Cat has been disabled.")
	}
}

func IsEnabled() bool {
	return atomic.LoadUint32(&isEnabled) > 0
}

func Shutdown() {
	scheduler.shutdown()
}

func DebugOn() {
	logger.logger.SetOutput(os.Stdout)
}
