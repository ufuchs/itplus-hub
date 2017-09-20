package app

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"hidrive.com/ufuchs/itplus/base/fcc"
	"hidrive.com/ufuchs/itplus/hub/device"
)

type (
	//
	MetaService struct {
		configDir    string
		metaFilename string
	}

	//
	DeviceConfigMeta struct {
		Hostname   string `yaml:"hostname"`
		Configname string `yaml:"configname"`
	}

	DeviceConfigMetas struct {
		ConfigMetas []DeviceConfigMeta `yaml:"configurations"`
	}
)

//
//
//
// func GetConfigname(hostname string, metas []DeviceConfigMeta) (string, error) {
// 	hostname = strings.ToLower(hostname)
// 	for _, m := range metas {
// 		if hostname == strings.ToLower(m.Hostname) {
// 			return m.Configname, nil
// 		}
// 	}
// 	s := fmt.Sprintf("Missing configuration for '%v'\n", hostname)
// 	return "", errors.New(s)
// }

const META_FILENAME = "meta.yml"
const CONFFILES_DIR = "configFiles"

//
// NewServices
//
func NewMetaService(configFilesDir string) (*MetaService, error) {

	if len(configFilesDir) == 0 {
		return nil, errors.New("Param 'ConfFilesDir' is empty")
	}

	_, err := os.Stat(configFilesDir)
	if err != nil {
		return nil, errors.New("'ConfFilesDir' doesn't exist")
	}

	filename := path.Join(configFilesDir, META_FILENAME)
	_, err = os.Stat(filename)
	if err != nil {
		msg := fmt.Sprintf("Missing file: '%v'", filename)
		return nil, errors.New(msg)
	}

	return &MetaService{
		configDir:    configFilesDir,
		metaFilename: filename,
	}, nil
}

// https://github.com/aerth/playwav

//
// RetreiveMeta
//
func (s *MetaService) RetreiveMeta() (*DeviceConfigMetas, error) {

	meta := &DeviceConfigMetas{}
	err := readYML(meta, s.metaFilename)
	return meta, err

}

//
//
//
func (s *MetaService) GetConfigname(hostname string, metas []DeviceConfigMeta) (string, error) {
	hostname = strings.ToLower(hostname)
	for _, m := range metas {
		if hostname == strings.ToLower(m.Hostname) {
			return m.Configname, nil
		}
	}
	return "", fmt.Errorf("    Missing configuration for '%v'\n", hostname)
}

//
//
//
func GetRuninfo(configFilesDir, configName string, offs int) (*fcc.RunInfo, error) {

	var (
		err error
		c   fcc.Configuration
		ri  *fcc.RunInfo
		dcs device.IConfigService
	)

	if dcs, err = device.NewConfigServiceEx(configFilesDir, configName, offs); err != nil {
		return nil, err
	}
	c = dcs.GetConfiguration()
	s := fmt.Sprintf("    Jeelink    : Configuration is { Name: %v, ModifiedOn: %v }", c.Name, c.ModifiedOn)
	fmt.Println(s)

	ri = fcc.NewRunInfoFromConfig(&c)

	return ri, err

}
