package jeelink

import (
	"testing"
)

//
//
//
func TestNewSerialConfig(t *testing.T) {

	_, err := NewSerialConfig()
	if err != nil {
		t.Error(err)
	}

}

//
//
//
func TestIsDeviceConnected(t *testing.T) {

	c, err := NewSerialConfig()
	if err != nil {
		t.Error(err)
	}

	err = IsDeviceConnected(c)
	if err != nil {
		t.Error(err)
	}

}
