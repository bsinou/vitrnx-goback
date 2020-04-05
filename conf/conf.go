package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var authURL string

// GetPort simply returns the main port to use.
func GetPort() string {
	port := viper.GetString(KeyPort)
	if port == "" {
		return "8888"
	}
	return port
}

// GetAllowedOrigins simply returns the correct value for CORS policy.
func GetAllowedOrigins() []string {
	orig := []string{viper.GetString(KeyPublicURL), viper.GetString(KeyBackendURL)}
	return orig
}

// GetKnownAdmins returns temporarily a list of admins.
func GetKnownAdmins() []string {
	ka := strings.Split(viper.GetString(KeyKnownAdmins), ", ")
	return ka
}

// GetAuthURL simply returns the correct value to get config.
func GetAuthURL() string {

	if len(authURL) > 1 {
		return authURL
	}

	authURL = fmt.Sprintf("%s/%s", viper.GetString(KeyGetAuthDN), viper.GetString(KeyGetAuthSuffix))
	return authURL
}

// GetJWKSURL returns the well known jwks retrieval URL for the configured Auth server.
func GetJWKSURL() string {
	return fmt.Sprintf("%s/%s", GetAuthURL(), JWKSWellKnownSuffix)
}

// GetDataFolderPath return the path to an existing data folder
func GetDataFolderPath() string {
	knownPaths := []string{
		// Deployed on Linux
		filepath.Join("/var", "lib", BaseName, VitrnxInstanceID, "data"),
		// In User Home on linux
		filepath.Join(os.Getenv("HOME"), ".config", BaseName, VitrnxInstanceID, "data"),
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
