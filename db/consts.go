package db

const (
	EMPTY = ""

	// debug config values
	DEBUG_CONFIG_DATABASE_NAME = "postgres"
	DEBUG_CONFIG_DATABASE_URL  = "host=localhost port=5432 " +
		"user=goilerplate password=goilerplate dbname=goilerplate"
	DEBUG_CONFIG_DATABASE_TEST = "host=localhost port=5432 " +
		"user=goilerplate password=goilerplate dbname=goilerplate_test"
	DEBUG_CONFIG_DATABASE_ADMIN = "host=localhost port=5432 user=postgres " +
		"password=postgres dbname=postgres"

	// config names
	CONFIG_DATABASE_NAME  = "database.name"
	CONFIG_DATABASE_URL   = "database.url"
	CONFIG_DATABASE_TEST  = "database.test"
	CONFIG_DATABASE_ADMIN = "database.admin"
)
