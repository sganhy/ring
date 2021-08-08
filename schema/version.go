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
	currentPatchVersion      int    = 4
	currentpreReleaseVersion string = "beta.2"

	// display version
	displayVersion string = "v%d.%d.%d"
)

//******************************
// getters and setters
//******************************

//******************************
// public methods
//******************************
func (ver *version) String() string {
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
