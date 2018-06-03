// Package conf provide easy configuration and 12-factor app compliant management of the VitrnX App
package conf

import (
	"fmt"
	"log"
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
	VitrnxVersion = "0.1.4-SNAPSHOT"
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
			return filepath.Join(p, fname), nil
		} else {
			fmt.Println(err.Error())
			// fmt.Printf("No file named %s found in %s: %v\n", fname, p, err.Error())
		}
	}
	return "", fmt.Errorf("could not find %s in known config folders", fname)
}

// GetKnownConfFolderPaths returns an array with the path of the folder
// where we might find the various configuration files
// Note that first found must be used.
func GetKnownConfFolderPaths() []string {
	deployed := filepath.Join("/etc", BaseName, VitrnxInstanceID)
	dev, err := filepath.Abs("conf")
	if err != nil {
		fmt.Println("Could not generate absolute path to code conf file: " + err.Error())
		dev = "conf" // Put a dummy value to avoid panic
	}

	if user, err := user.Current(); err == nil {
		return []string{
			// Deployed on Linux
			deployed,
			// In User Home on linux
			filepath.Join(user.HomeDir, ".config", BaseName, VitrnxInstanceID),
			// In source code in dev mode
			dev,
		}
	}

	return []string{
		// Deployed on Linux
		deployed,
		// In source code in dev mode
		dev,
	}
}

// GetDataFolderPath return the path to an existing data folder
func GetDataFolderPath() string {
	knownPaths := []string{
		// Deployed on Linux
		filepath.Join("/var", "lib", BaseName, VitrnxInstanceID, "data"),
		// In User Home on linux
		// filepath.Join("$HOME", ".config", BaseName, VitrnxInstanceID, "data"),
		// In source code in dev mode
		filepath.Join(".", "data"),
	}

	for _, path := range knownPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	log.Fatalf("no existing data folder found, please configure your environment")
	return ""
}
