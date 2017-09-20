package device

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type (
	//
	Timeline struct {
		configDir string
	}

	ITimeline interface {
		GetTimelines() (Configurations, error)
		GetMembers() (Members, error)
		GetTimelineBy(string) ([]string, error)
		GetSingleTimelineBy(string, int) (string, error)
	}

	//
	Configurations map[string][]string
	//
	Member struct {
		Name  string `yaml:"name" json:"name"`
		Count int    `yaml:"count" json:"count"`
	}
	Members []Member
)

func (slice Members) Len() int           { return len(slice) }
func (slice Members) Less(i, j int) bool { return slice[i].Name < slice[j].Name }
func (slice Members) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

//
//
//
func NewTimeline(configDir string) ITimeline {
	return &Timeline{
		configDir: configDir,
	}
}

//
// fetchConfigFilenames returns all filenames of the 'profile' dir
//
func fetchConfigFilenames(configDir string) ([]string, error) {

	filenames := make([]string, 0)

	if len(configDir) == 0 {
		return nil, ErrParamConfigDirIsEmpty()
	}

	fileInfos, err := ioutil.ReadDir(configDir)
	if err != nil {
		return filenames, err
	}

	for _, fi := range fileInfos {

		if !fi.IsDir() {

			var name = fi.Name()

			var pos = strings.Index(name, "-")
			if pos < 0 {
				continue
			}

			var ext = filepath.Ext(name)
			if ext != ".yml" {
				continue
			}

			name = name[0 : len(name)-len(ext)]

			filenames = append(filenames, name)
		}
	}

	if len(filenames) == 0 {
		return filenames, ErrNoConfigfilesFound(configDir)
	}

	sort.Strings(filenames)

	return filenames, nil
}

//
//
//
func sortFilenamesToConfigurations(filenames []string) (configurations Configurations) {

	configurations = make(Configurations)
	for _, filename := range filenames {

		var pos = strings.Index(filename, "-")
		var key = filename[0:pos]
		var value = configurations[key]

		configurations[key] = append(value, filename)
	}

	return

}

//
// GetConfigNames returns the basenames of the profiles dir,
// e.g. ['home', 'aleta']
//
func retrieveAll(configDir string) (configurations Configurations, err error) {

	var filenames []string

	if filenames, err = fetchConfigFilenames(configDir); err != nil {
		return nil, err
	}

	configurations = sortFilenamesToConfigurations(filenames)

	return configurations, nil

}

//
// GetTimeline return a reverse order of the filenames of a given configuration
//
func (p *Timeline) GetTimelines() (timeline Configurations, err error) {

	var configurations Configurations

	if configurations, err = retrieveAll(p.configDir); err != nil {
		return nil, err
	}

	timeline = make(Configurations)

	for k, v := range configurations {

		var files sort.StringSlice = v
		files.Sort()
		sort.Sort(sort.Reverse(files[:]))

		timeline[k] = files

	}

	return
}

//
//
//
func (p *Timeline) GetTimelineBy(name string) ([]string, error) {

	if name == "" {
		return nil, ErrMissingTimelineName()
	}

	timelines, err := p.GetTimelines()
	if err != nil {
		return nil, err
	}

	result, ok := timelines[name]
	if !ok {
		return nil, ErrTimelineDoesntExist(name)
	}

	return result, nil
}

//
//
//
func (p *Timeline) GetSingleTimelineBy(name string, offs int) (string, error) {

	if offs < 0 {
		offs = 0
	}

	timeline, err := p.GetTimelineBy(name)
	if err != nil {
		return "", err
	}

	if offs > len(timeline)-1 {
		offs = len(timeline) - 1
	}

	return timeline[offs], nil

}

//
// GetConfigNames returns the basenames of the profiles dir,
// e.g. 'home', 'company'
//
func (p *Timeline) GetMembers() (members Members, err error) {

	var c Configurations

	if c, err = retrieveAll(p.configDir); err != nil {
		return nil, err
	}

	members = make(Members, 0)
	for k, v := range c {
		members = append(members, Member{k, len(v)})
	}

	sort.Sort(members)

	return

}
