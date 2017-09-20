package app

import (
	"fmt"
	"path"
	"testing"
)

var conffiledir = path.Join("..", "configFiles")

//
//
//
func TestMetaService(t *testing.T) {

	_, err := NewMetaService(conffiledir)
	if err != nil {
		t.Error(err)
	}

}

//
//
//
func TestMetaServiceRetrieve(t *testing.T) {

	s, err := NewMetaService(conffiledir)
	if err != nil {
		t.Error(err)
		return
	}

	var meta *DeviceConfigMetas
	meta, err = s.RetreiveMeta()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(meta)

	for _, e := range meta.ConfigMetas {
		fmt.Println(e)
	}

	_, err = GetConfigname("salata", meta.ConfigMetas)
	if err != nil {
		t.Error(err)
		return
	}

}
