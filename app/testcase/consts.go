package testcase

const (
	EMPTY = ""

	// testcase
	DROP_DATABASE_STATEMENT   = `DROP DATABASE %v;`
	CREATE_DATABASE_STATEMENT = `CREATE DATABASE %v;`
	SET_GRANT                 = `grant ALL privileges on database "%v" to "%v"`
	COOKIE                    = "access"
	SEND_BAD_JSON             = "when send bad json"
	BAD_JSON                  = `{"badjson}`
)
