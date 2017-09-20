package device

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"hidrive.com/ufuchs/itplus/base/fcc"
	"hidrive.com/ufuchs/itplus/base/oops"
)

var testFilenames = []string{
	"aleta-2017-05-15T202418Z",
	"aleta-2017-05-17T155848Z",
	"aleta-2017-05-17T175852Z",
	"aleta-2017-05-17T175924Z",
	"aleta-2017-05-17T182418Z",
	"home-2017-05-18T160934Z",
}

var testConfigurations = Configurations{
	"aleta": {
		"aleta-2017-05-15T202418Z",
		"aleta-2017-05-17T155848Z",
		"aleta-2017-05-17T175852Z",
		"aleta-2017-05-17T175924Z",
		"aleta-2017-05-17T182418Z"},
	"home": {
		"home-2017-05-18T160934Z"},
}

var testMembers = Members{
	{Name: "aleta", Count: 5},
	{Name: "home", Count: 1},
}

var aletaReversed = []string{
	"aleta-2017-05-17T182418Z",
	"aleta-2017-05-17T175924Z",
	"aleta-2017-05-17T175852Z",
	"aleta-2017-05-17T155848Z",
	"aleta-2017-05-15T202418Z",
}

// /////////////////////////////////////////////////////////////////////////////
// helper
// /////////////////////////////////////////////////////////////////////////////

func createTestFiles(profileDir string) {

	for _, filename := range testFilenames {

		filename = path.Join(profileDir, filename+".yml")

		err := ioutil.WriteFile(filename, []byte{}, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

}

// /////////////////////////////////////////////////////////////////////////////
// tests
// /////////////////////////////////////////////////////////////////////////////

//
//
//
func TestFetchConfigFilenames(t *testing.T) {

	if _, err := fetchConfigFilenames(""); err == nil {
		t.Errorf("FetchConfigFilenames(): should throw ==> '%v'", "config dir name is empty")
		return
	}

	var noneexistingDir = "cafebabe"
	if _, err := fetchConfigFilenames(noneexistingDir); err == nil {
		t.Errorf("FetchConfigFilenames(): should throw ==> '%v'", "open cafebabe: no such file or directory")
		return
	}

	var emptyDir = path.Join(testConfigDir, "cafababe")
	os.Mkdir(emptyDir, 0644)
	if _, err := fetchConfigFilenames(emptyDir); err == nil {
		t.Errorf("FetchConfigFilenames(): should throw ==> '%v'", "Your config dir seems to be empty")
		return
	}
	os.Remove(emptyDir)

	//////////////////////////////////////////////////////////////////////////////

	expected := testFilenames

	createTestFiles(testConfigDir)

	actual, err := fetchConfigFilenames(testConfigDir)
	if err != nil {
		t.Errorf("FetchConfigFilenames(): err => '%v'", err)
		return
	}

	fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("GetFilenames(): expected '%v', actual '%v'", expected[i], actual[i])
		}
	}

}

//
//
//
func TestSortFilenamesToConfigurations(t *testing.T) {

	actual := sortFilenamesToConfigurations(testFilenames)

	for ak, av := range actual {

		ev := testConfigurations[ak]

		for i := range av {

			if ev[i] != av[i] {
				t.Errorf("SortFilenamesToConfigurations(): expected '%v', actual '%v'", ev[i], av[i])
				break
			}

		}

	}

}

//
//
//
func TestRetrieveAll(t *testing.T) {

	createTestFiles(testConfigDir)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	if _, err := retrieveAll("testConfigDir"); err == nil {
		t.Errorf("getConfigurations(): should throw ==> '%v'", "open testConfigDir: no such file or directory")
	}

	if _, err := retrieveAll(testConfigDir); err != nil {
		t.Errorf("getConfigurations(): '%v'", err)
	}

}

//
//
//
func TestGetMembers(t *testing.T) {

	p := NewTimeline("testConfigDir")
	if _, err := p.GetMembers(); err == nil {
		t.Errorf("GetMembers(): should throw ==> '%v'", "open testConfigDir: no such file or directory")
	}

	createTestFiles(testConfigDir)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	p = NewTimeline(testConfigDir)

	actual, _ := p.GetMembers()

	for i := range actual {
		if testMembers[i] != actual[i] {
			t.Errorf("GetMembers(): expected '%v', actual '%v'", testMembers[i], actual[i])
			break
		}
	}

	// check the members
	for _, c := range actual {
		switch c.Name {
		case "aleta":
			if c.Count != 5 {
				t.Error("GetMembers(): 'count' doesn't equals 5 for " + c.Name)
			}
		case "home":
			if c.Count != 1 {
				t.Error("GetMembers(): 'count' doesn't equals 1 for " + c.Name)
			}

		default:
			t.Error("GetMembers(): missing configuration for " + c.Name)
		}

	}

}

//
//
//
func TestGetTimelines(t *testing.T) {

	createTestFiles(testConfigDir)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	p := NewTimeline("testConfigDir")
	if _, err := p.GetTimelines(); err == nil {
		t.Errorf("GetTimelines() should throw ==> '%v'", "open testConfigDir: no such file or directory")
	}

	p = NewTimeline(testConfigDir)
	_, err := p.GetTimelines()
	if err != nil {
		t.Error(err)
	}

}

//
//
//
func TestGetTimelinesBy(t *testing.T) {

	p := NewTimeline("testConfigDir")
	if _, err := p.GetTimelineBy("aleta"); err == nil {
		t.Errorf("GetTimelineBy() should throw ==> '%v'", "open testConfigDir: no such file or directory")
		return
	}

	createTestFiles(testConfigDir)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	p = NewTimeline(testConfigDir)
	if _, err := p.GetTimelineBy(""); err == nil {
		t.Errorf("GetTimelineBy() should throw ==> '%v'", "Parameter 'name' is empty")
		return
	}

	p = NewTimeline(testConfigDir)
	_, err := p.GetTimelineBy("a")
	if err == nil {
		t.Errorf("GetTimelineBy() : err => '%v'", err)
		return
	}

	p = NewTimeline(testConfigDir)

	var timeline []string
	timeline, err = p.GetTimelineBy("aleta")
	if err != nil {
		t.Errorf("GetTimelineBy() : err => '%v'", err)
		return
	}

	for _, filename := range timeline {

		err := fileExists(testConfigDir, filename+".yml")
		if err != nil {
			t.Error(err)
		}

	}

}

func fileExists(configDir, filename string) error {

	filename = path.Join(configDir, filename)
	_, err := os.Stat(filename)

	return err

}

//
//
//
func TestGetSingleTimelineBy(t *testing.T) {

	createTestFiles(testConfigDir)
	defer fcc.RemoveFiles(TestMemberNames, testConfigDir, t)

	p := NewTimeline(testConfigDir)

	// offeset less than 0
	filename, err := p.GetSingleTimelineBy("aleta", -1)
	if err != nil {
		t.Error(err)
	}

	// offeset less than 0
	filename, err = p.GetSingleTimelineBy("aleta", 20)
	if err != nil {
		t.Error(err)
	}

	filename, err = p.GetSingleTimelineBy("aleta__", 1)
	if err == nil {
		s := fmt.Sprintf(getErrorDescription(oops.EC_TIMELINE_FOR_XY_DOESNT_EXIST), "aleta__")
		t.Errorf("GetSingleTimelineBy() should throw ==> %v", s)
	}

	err = fileExists(testConfigDir, filename)
	if err != nil {
		t.Error(err)
	}

}
