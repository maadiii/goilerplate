package controllers

const (
	EMPTY = ""
	SPACE = " "
	COMMA = ","
	COLON = ":"

	// methods
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	PATCH  = "PATCH"

	// routes
	ADMIN_ACCOUNTS        = "/admin/accounts"
	ADMIN_ACCOUNTS_CREATE = ADMIN_ACCOUNTS + "/create"
	ADMIN_GROUPS          = "/admin/groups"
	ADMIN_GROUPS_CREATE   = ADMIN_GROUPS + "/create"
	ADMIN_GROUPS_EDIT     = ADMIN_GROUPS + "/edit/:id"
	APIV1                 = "/apiv1"

	// root
	HTTPS_PORT              = "ports.https"
	HTTP_PORT               = "ports.http"
	DOMAIN_NAME             = "domain-name"
	PROXY_COUNT             = "proxy-count"
	DURATION                = "duration"
	STATUS_CODE             = "status_code"
	REMOTE                  = "remote"
	CONTENT_TYPE            = "Content-Type"
	JSON_CONTENT_TYPE       = "application/json"
	TEXT_PLAIN_CONTENT_TYPE = "text/plain; charset=utf-8"
	TEXT_HTML_CONTENT_TYPE  = "text/html"
	NO_SNIFF                = "nosniff"
	X_CONTENT_TYPE_OPTIONS  = "X-Content-Type-Options"
	X_FORWARDED_FOR         = "X-Forwarded-For"
	ACCESS_COOKIE           = "access"
	REFRESH_COOKIE          = "refresh"

	// error messages
	HTTP_ERROR  = "ports.http is not set in config file"
	HTTPS_ERROR = "ports.https is not set in config file"
)
