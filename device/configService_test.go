package device

import (
	"os"
	"strings"
	"testing"

	"hidrive.com/ufuchs/itplus/base/fcc"
)

var testConfigDir = os.TempDir()
var TestMemberNames = []string{"aleta", "home"}

// /////////////////////////////////////////////////////////////////////////////
// helper
// /////////////////////////////////////////////////////////////////////////////

func createFiles(profileDir string, t *testing.T) {

	dao := fcc.NewConfigDAO(profileDir)

	////
	// test, first create the files,
	//       second read them
	////

	for k, v := range fcc.AletaYML {

		config, err := fcc.UnmarshalDeviceConfigYML(v)
		if err != nil {
			t.Error("unmarshalDeviceConfigYML():", err)
		}

		if err = dao.WriteConfig(config, k); err != nil {
			t.Error("dao.writeConfig():", err)
		}

	}

}

// /////////////////////////////////////////////////////////////////////////////
// tests
// /////////////////////////////////////////////////////////////////////////////

//
//
//
func TestITimeline(t *testing.T) {

	createFiles(testConfigDir, t)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	dcs := NewConfigService(testConfigDir)

	if _, err := dcs.GetMembers(); err != nil {
		t.Error(err)
	}

	if _, err := dcs.GetTimelines(); err != nil {
		t.Error(err)
	}

	if _, err := dcs.GetTimelineBy("aleta"); err != nil {
		t.Error(err)
	}

	if _, err := dcs.GetSingleTimelineBy("aleta", 0); err != nil {
		t.Error(err)
	}

}

//
//
//
func TestIConfiguration_Retrieve(t *testing.T) {

	// setup

	createFiles(testConfigDir, t)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	dcs := NewConfigService(testConfigDir)

	configFile, err := dcs.GetSingleTimelineBy("aleta", 0)
	if err != nil {
		t.Error(err)
	}

	// action

	err = dcs.Retrieve("configFile")
	if err == nil {
		t.Error("Retrieve() should throw => 'no such file or directory'")
	}

	err = dcs.Retrieve(configFile)
	if err != nil {
		t.Error(err)
	}

	// test

	var pos = strings.Index(configFile, "-")
	var timestamp = configFile[pos+1:]

	if dcs.GetConfiguration().ModifiedOn != timestamp {
		t.Error("args")
	}

}

//
//
//
func TestIConfiguration_RemoveDevice(t *testing.T) {

	// setup

	createFiles(testConfigDir, t)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	dcs, err := NewConfigServiceEx(testConfigDir, "aleta", 0)
	if err != nil {
		t.Error(err)
	}

	// setup

	var id = 9
	err = dcs.RemoveDevice(id)
	if err == nil {
		t.Error("RemoveDevice(): should throw => ITPLUS-0302: Device doesn't exists 'Id = 9'")
	}

	id = 1
	err = dcs.RemoveDevice(id)
	if err != nil {
		t.Errorf("RemoveDevice(): %v", err)
	}

	// check

	// get latest configuration, it's the new one
	var configFile string
	configFile, err = dcs.GetSingleTimelineBy("aleta", 0)
	if err != nil {
		t.Error(err)
	}

	err = dcs.Retrieve(configFile)
	if err != nil {
		t.Error(err)
	}

	_, err = dcs.GetDevice(id)
	if err == nil {
		t.Error("RemoveDevice(): after re-retrieve device 'Id = 1' -> it is still listed!")
	}

}

//
//
//
func TestIConfiguration_UpdateDeviceByJSON(t *testing.T) {

	// setup

	createFiles(testConfigDir, t)

	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	dcs, err := NewConfigServiceEx(testConfigDir, "aleta", 0)
	if err != nil {
		t.Error(err)
	}

	// test

	var sameOne = `{
			"id": 1,
			"num": 200,
			"addr": 7,
			"type": "TX29TDH-IT",
			"alias": "R 2.7 - Schreibtisch",
			"absent": false,
			"locked": false,
			"lon": 19.31,
			"lat": 19.32,
			"alt": 19.33
	}`

	err = dcs.UpdateDeviceByJSON([]byte(sameOne))
	if err != nil {
		t.Errorf("UpdateDeviceByJSON(): %v", err)
	}

	// check

	// get latest configuration, it's the new one
	var configFile string
	configFile, err = dcs.GetSingleTimelineBy("aleta", 0)
	if err != nil {
		t.Error(err)
	}

	// fmt.Println(configFile)

	err = dcs.Retrieve(configFile)
	if err != nil {
		t.Error(err)
	}

	// fmt.Println(dcs.GetConfiguration())

	var device fcc.Device
	device, err = dcs.GetDevice(1)
	if device.Num != 200 {
		t.Error("UpdateDeviceByJSON(): expected 'Num = 200'")
	}

}

//
//
//
func TestIConfiguration_AddDeviceByJSON(t *testing.T) {

	// setup

	createFiles(testConfigDir, t)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	dcs, err := NewConfigServiceEx(testConfigDir, "aleta", 0)
	if err != nil {
		t.Error(err)
	}

	var newOne = `{
			"id": 1,
			"num": 200,
			"addr": 123,
			"type": "TX29TDH-IT",
			"alias": "R 2.7 - Schreibtisch",
			"absent": false,
			"locked": false,
			"lon": 19.91,
			"lat": 19.92,
			"alt": 19.93
	}`

	var id int
	id, err = dcs.AddDeviceByJSON([]byte(newOne))
	if err != nil {
		t.Errorf("AddDeviceByJSON(): %v", err)
	}

	// check

	// get latest configuration, it's the new one
	var configFile string
	configFile, err = dcs.GetSingleTimelineBy("aleta", 0)
	if err != nil {
		t.Error(err)
	}

	err = dcs.Retrieve(configFile)
	if err != nil {
		t.Error(err)
	}
	var d fcc.Device
	d, err = dcs.GetDevice(id)
	if d.Num != 200 && d.Addr != 123 {
		t.Error("AddDeviceByJSON(): expected 'Num = 200' and 'Addr = 123'")
	}

}
