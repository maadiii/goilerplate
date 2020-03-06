package cli

const (
	EMPTY = ""
	COLON = ":"

	// database
	DB_USE         = "db"
	DB_SHORT       = "Database utilities"
	MIGRATE_USE    = "migrate"
	MIGRATE_SHORT  = "Migrates database"
	DROP_USE       = "drop"
	DROP_SHORT     = "Drop database table"
	BASEDATA_USE   = "basedata"
	BASEDATA_SHORT = "Insert base data to database"

	//encrytption
	CRYPTO_USE           = "crypto"
	CRYPTO_SHORT         = "Cryptography utilities"
	SYMMETRIC_USE        = "sym-key"
	SYMMETRIC_SHORT      = "Generate symmetric secret key"
	JWT_USE              = "jwt"
	JWT_SHORT            = "Create json web token"
	LENGTH_FLAG          = "length"
	LENGTH_FLAG_SHORT    = "l"
	LENGTH_FLAG_MESSAGE  = "Length of bits"
	SECONDS_FLAG         = "seconds"
	SECONDS_FLAG_SHORT   = "s"
	SECONDS_FLAG_MESSAGE = "Max age of JWT"

	// root
	VERBOSE_FLAG         = "verbose"
	VERBOSE_FLAG_SHORT   = "v"
	VERBOSE_FLAG_MESSAGE = "make output verbose"

	// version
	VERSION_USE   = "version"
	VERSION_SHORT = "Print the version number"

	// serve
	HTTPS_SCHEME        = "https://"
	STATIC_PATH         = "/static/"
	CONTENT_TYPE        = "Content-Type"
	JSON_CONTENT_TYPE   = "application/json"
	SERVING_LOG_INFO    = "Serving at http://127.0.0.1:%d"
	SERVE_USE           = "serve"
	SERVE_SHORT         = "Serves the application"
	SHUTING_DOWN_SIGNAL = "Signal caught. Shutting down..."
	CONFIG_FLAG         = "config"
	CONFIG_FLAG_SHORT   = "c"
	CONFIG_FLAG_USAGE   = "set config file in .yml" +
		"(default is ./app/config.go -> var debugConfig)"
)
