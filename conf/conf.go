package conf

const (
	// EnvDev is the development environment
	EnvDev = "dev"
	// EnvTest is the test and staging environment
	EnvTest = "test"
	// EnvProd is the production environment
	EnvProd = "production"
)

var (
	// The 3 below vars are initialized by the go linker directly
	// in the resulting binary when doing 'make main'

	// VitrnxInstanceKey exposes the unique name of the current instance,
	// it is used among others in various file paths when in prod environment
	VitrnxInstanceKey = "vitrnx"
	// VitrnxVersion exposes current version of the backend
	VitrnxVersion = "0.1.0"
	// BuildTimestamp exposes running app build time stamp
	BuildTimestamp = ""
	// BuildRevision exposes the git id that was at master origin head at the time of building
	BuildRevision = ""
	// Env exposes current environment type (dev, prod...)
	Env = EnvDev
)
