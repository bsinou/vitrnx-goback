package conf

import (
	"github.com/spf13/viper"
)

// GetAllowedOrigins simply returns the correct value for CORS policy.
func GetAllowedOrigins() string {
	return viper.GetString(KeyPublicURL)
}
