//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Geduld ist eine gute Eigenschaft. Aber nicht, wenn es um die Beseitigung ]
// [ von Missständen geht.                                                    ]
// [                                                      -Margaret Thatcher- ]

package app

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const (
	FILENAME = "app.yml"
)

type (
	ConfigService struct {
		LastErr error
	}
)

var (
	BaseDir        string
	DataPort       int
	DiscoveryPort  int
	ServicePort    int
	DefaultConfig  string
	ConfigFilesDir string
)

//
//
//
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

//
//
//
func (d *ConfigService) RetrieveAll() *ConfigService {

	type Config struct {
		DataPort       int    `yaml:"dataport"`
		DiscoveryPort  int    `yaml:"discoveryport"`
		ServicePort    int    `yaml:"serviceport"`
		ConfigFilesDir string `yaml:"configfilesdir"`
	}

	if d.LastErr != nil {
		return d
	}

	filename := path.Join(BaseDir, FILENAME)

	config := &Config{}
	d.LastErr = readYML(config, filename)

	DataPort = config.DataPort
	DiscoveryPort = config.DiscoveryPort
	ServicePort = config.ServicePort

	if len(config.ConfigFilesDir) != 0 {
		ConfigFilesDir = config.ConfigFilesDir
	} else {
		ConfigFilesDir = path.Join(BaseDir, CONFFILES_DIR)
	}

	return d
}

//
//
//
func readYML(v interface{}, filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(raw, v)
}
