package device

import (
	"ufuchs/itplus/base/fcc"
)

type (
	DeviceConfigService struct {
		configDir     string
		timeline      ITimeline
		configuration fcc.Configuration
	}

	IConfigService interface {
		ITimeline
		fcc.IConfiguration
		Retrieve(string) error
		GetConfiguration() fcc.Configuration
	}
)

//
//
//
func NewConfigService(configDir string) IConfigService {
	return &DeviceConfigService{
		configDir:     configDir,
		timeline:      NewTimeline(configDir),
		configuration: fcc.Configuration{},
	}
}

//
//
//
func NewConfigServiceEx(configDir, memberName string, offs int) (IConfigService, error) {

	dcs := NewConfigService(configDir)

	configFile, err := dcs.GetSingleTimelineBy(memberName, offs)
	if err != nil {
		return nil, err
	}

	return dcs, dcs.Retrieve(configFile)

}

// ITimeline ///////////////////////////////////////////////////////////////////

//
//
//
func (dcs *DeviceConfigService) GetMembers() (Members, error) {
	return dcs.timeline.GetMembers()
}

//
//
//
func (dcs *DeviceConfigService) GetTimelineBy(memberName string) ([]string, error) {
	return dcs.timeline.GetTimelineBy(memberName)
}

//
//
//
func (dcs *DeviceConfigService) GetSingleTimelineBy(memberName string, offs int) (string, error) {
	return dcs.timeline.GetSingleTimelineBy(memberName, offs)
}

//
//
//
func (dcs *DeviceConfigService) GetTimelines() (timeline Configurations, err error) {
	return dcs.timeline.GetTimelines()
}

// IConfiguration //////////////////////////////////////////////////////////////

//
//
//
func (dcs *DeviceConfigService) AddDeviceByJSON(newOne []byte) (id int, err error) {

	if id, err = dcs.configuration.AddDeviceByJSON(newOne); err != nil {
		return
	}

	dao := fcc.NewConfigDAO(dcs.configDir)
	if err = dao.WriteConfig(&dcs.configuration, ""); err != nil {
		id = -1
	}

	return id, err
}

//
//
//
func (dcs *DeviceConfigService) GetDevice(id int) (fcc.Device, error) {
	return dcs.configuration.GetDevice(id)
}

//
//
//
func (dcs *DeviceConfigService) UpdateDeviceByJSON(updateOne []byte) error {

	if err := dcs.configuration.UpdateDeviceByJSON(updateOne); err != nil {
		return err
	}

	return fcc.NewConfigDAO(dcs.configDir).
		WriteConfig(&dcs.configuration, "")

}

//
//
//
func (dcs *DeviceConfigService) RemoveDevice(id int) error {

	if err := dcs.configuration.RemoveDevice(id); err != nil {
		return err
	}

	return fcc.NewConfigDAO(dcs.configDir).
		WriteConfig(&dcs.configuration, "")

}

//
//
//
func (dcs *DeviceConfigService) GetConfiguration() fcc.Configuration {
	return dcs.configuration
}

//
//
//
func (dcs *DeviceConfigService) Retrieve(configName string) error {

	dao := fcc.NewConfigDAO(dcs.configDir)

	var err error
	dcs.configuration, err = dao.Retrieve(configName)

	return err

}
