package conf

/* MAIN CONSTANTS */

const (
	// BaseName holds a constant to the vitrnx name in case it should ever be changed.
	BaseName = "vitrnx"
)

/* KEYS (to retrieve values via Viper) */
const (
	KeyEnvType        = "env"
	KeyAdminEmail     = "auth.adminEmail"
	KeyAnonymousEmail = "auth.anonymousEmail"
	KeyKnownRoles     = "auth.knownRoles"

	KeyPublicURL = "publicUrl"
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
