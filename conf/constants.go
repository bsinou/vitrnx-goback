package conf

/* MAIN CONSTANTS */

const (
	BaseName = "vitrnx"
)

/* KEYS (to retrieve values via Viper) */
const (
	KeyEnvType = "env"

	KeyAdminUsers = "users.admins"
)

/* KNOWN VALUES */

const (
	// EnvDev is the development environment
	EnvDev = "dev"
	// EnvTest is the test environment
	EnvTest = "test"
	// EnvStaging is the pre-production environment
	EnvStaging = "staging"
	// EnvProd is the production environment
	EnvProd = "prod"
)
