package conf

/* MAIN CONSTANTS */

const (
	// BaseName holds a constant to the vitrnx name in case it should ever be changed.
	BaseName = "vitrnx"
)

/* KEYS (to retrieve values via Viper) */
const (
	KeyEnvType = "env"
	KeyPort    = "port"

	KeyAdminEmail     = "auth.adminEmail"
	KeyAnonymousEmail = "auth.anonymousEmail"
	KeyKnownRoles     = "auth.knownRoles"
	KeyKnownAdmins    = "auth.knownAdmins"

	KeyPublicURL  = "publicUrl"
	KeyBackendURL = "backendUrl"

	KeyGetAuthDN     = "authDN"
	KeyGetAuthSuffix = "authSuffix"
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

	// OIDC Entrypoints
	JWKSWellKnownSuffix = ".well-known/jwks.json"
)
