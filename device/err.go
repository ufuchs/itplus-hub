package device

import (
	"fmt"

	"ufuchs/itplus/base/oops"
)

var EC_Messages = map[int]string{
	//

	oops.EC_PARAM_CONFIG_DIR_EMPTY:       "param 'configDir' is empty",
	oops.EC_NO_CONFIGFILES_FOUND:         "the configDir '%v' doesn't contain any YML files",
	oops.EC_MISSING_TIMELINE_NAME:        "name of the timeline is empty",
	oops.EC_TIMELINE_FOR_XY_DOESNT_EXIST: "a timeline for '%v' doesn't exists",
}

var ErrTimelineDoesntExist = func(configName string) error {
	errno := oops.EC_TIMELINE_FOR_XY_DOESNT_EXIST
	msg := fmt.Sprintf(getErrorDescription(errno), configName)
	return &oops.Err{errno, msg}
}

var ErrMissingTimelineName = func() error {
	errno := oops.EC_MISSING_TIMELINE_NAME
	msg := getErrorDescription(errno)
	return &oops.Err{errno, msg}
}

var ErrParamConfigDirIsEmpty = func() error {
	errno := oops.EC_PARAM_CONFIG_DIR_EMPTY
	msg := getErrorDescription(errno)
	return &oops.Err{errno, msg}
}

var ErrNoConfigfilesFound = func(dirname string) error {
	errno := oops.EC_NO_CONFIGFILES_FOUND
	msg := fmt.Sprintf(getErrorDescription(errno), dirname)
	return &oops.Err{errno, msg}
}

// GetDescription returns the corresponding verbal description of 'code'
func getErrorDescription(errno int) string {
	desc, ok := EC_Messages[errno]
	if !ok {
		desc = oops.UNKNOWN_ERRNO
	}
	return desc
}
