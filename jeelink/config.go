package jeelink

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tarm/serial"
)

//
//
//
func NewSerialConfig() (*serial.Config, error) {

	var device string

	if device = os.Getenv(JEELINK_PORT); device == "" {
		return nil, ErrMissingJEELINKEnvironment
	}

	fmt.Println("==> JEELINK device from ENV is:", device)

	return &serial.Config{Name: device, Baud: BAUD}, nil

}

//
//
//
func IsDeviceConnected(config *serial.Config) error {

	switch runtime.GOOS {
	case "windows":
		port, err := serial.OpenPort(config)
		if err != nil {
			return ErrMissingJEELINK
		}
		if err = port.Close(); err != nil {
			return err
		}
	default:
		if _, err := os.Stat(config.Name); err != nil {
			return ErrMissingJEELINK
		}

	}

	return nil

}
