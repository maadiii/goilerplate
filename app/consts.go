package app

const (
	EMPTY       = ""
	NAME        = "goilerplate"
	DESCRIPTION = "Goilerplate web application"
	VERSION     = "0.1.0dev"

	// config
	SECRET_KEY              = "secret-key"
	BLOCK_SECRET_KEY        = "block-secret-key"
	JWT_SECRET              = "jwt.secret"
	JWT_ALGORITHM           = "jwt.algorithm"
	JWT_MAXAGE              = "jwt.max_age"
	JWT_HTTPONLY            = "jwt.httponly"
	REFRESH_TOKEN_SECRET    = "jwt.refresh_token.secret"
	REFRESH_TOKEN_ALGORITHM = "jwt.refresh_token.algorithm"
	REFRESH_TOKEN_MAXAGE    = "jwt.refresh_token.max_age"
	REFRESH_TOKEN_SECURE    = "jwt.refresh_token.secure"
	REFRESH_TOKEN_HTTPONLY  = "jwt.refresh_token.httponly"
	REFRESH_TOKEN_PATH      = "jwt.refresh_token.path"
	STATIC                  = "static"
	CRT                     = "tls.crt"
	KEY                     = "tls.key"
	HS256                   = "HS256"
	HS384                   = "HS384"
	HS512                   = "HS512"

	// principles
	DIGITS   = "0123456789"
	SPECIALS = "~=+%^*/()[]{}/!@#$?|<>"
	ALL      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" +
		DIGITS + SPECIALS
	ALL_UNSPECIALS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" +
		DIGITS

	// Error messages
	UNAUTHORIZED = "Unauthorized"
	FORBIDDEN    = "Foribidden"
)
