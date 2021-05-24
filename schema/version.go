package schema

import (
	"fmt"
)

type version struct {
	majorVersion      int
	minorVersion      int
	patchVersion      int
	preReleaseVersion string
}

const (
	currentMajorVersion      int    = 1
	currentMinorVersion      int    = 0
	currentPatchVersion      int    = 3
	currentpreReleaseVersion string = "beta.2"

	// display version
	displayVersion string = "v%d.%d.%d"
)

//******************************
// getters
//******************************

//******************************
// public methods
//******************************
func (ver *version) GetCurrentVersion() string {
	ver.getCurrentVersion()
	var result = fmt.Sprintf(displayVersion, ver.majorVersion, ver.minorVersion, ver.patchVersion)
	if currentpreReleaseVersion != "" {
		result += "-" + ver.preReleaseVersion
	}
	return result
}

//******************************
// private methods
//******************************
func (ver *version) getCurrentVersion() {
	ver.majorVersion = currentMajorVersion
	ver.minorVersion = currentMinorVersion
	ver.patchVersion = currentPatchVersion
	ver.preReleaseVersion = currentpreReleaseVersion
}
