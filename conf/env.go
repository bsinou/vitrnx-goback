// Package conf provide easy configuration and 12-factor app compliant management of the VitrnX App
package conf

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var (
	// VitrnxInstanceID exposes the unique name of the current instance,
	// it is used among others in various file paths when in prod environment
	VitrnxInstanceID string

	// The 3 below vars are initialized by the go linker directly
	// in the resulting binary when doing 'make main'

	// VitrnxVersion exposes current version of the backend
	VitrnxVersion = "0.0.1"
	// BuildTimestamp exposes running app build time stamp
	BuildTimestamp = ""
	// BuildRevision exposes the git id that was at master origin head at the time of building
	BuildRevision = ""

	// Env exposes current environment type (dev, prod...)
	Env = EnvDev
)

func GetConfigFile(fname string) (string, error) {
	for _, p := range GetKnownConfFolderPaths() {
		cp := filepath.Join(p, fname)
		// fmt.Printf("Looking for %s, checking %s: [%s]\n", fname, p, cp)
		if _, err := os.Stat(cp); err == nil {
			fmt.Println("Using configuration file at", cp)
			return filepath.Join(p, fname), nil
		}
		// fmt.Println(err.Error())
	}
	return "", fmt.Errorf("could not find %s in known config folders", fname)
}

// GetKnownConfFolderPaths returns an array with the path of the folder
// where we might find the various configuration files
// Note that first found must be used.
func GetKnownConfFolderPaths() (paths []string) {

	// Deployed on Linux
	paths = append(paths, filepath.Join("/etc", BaseName, VitrnxInstanceID))

	// In User Home on linux
	if user, err := user.Current(); err == nil {
		paths = append(paths, filepath.Join(user.HomeDir, ".config", BaseName, VitrnxInstanceID))
	}

	// In source code in dev mode
	dev, err := filepath.Abs("conf")
	if err != nil {
		fmt.Println("Could not generate absolute path to code conf file: " + err.Error())
		dev = "conf" // Put a dummy value to avoid panic
	}

	return append(paths, dev)
}
