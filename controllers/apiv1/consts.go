package apiv1

const (
	APPLICATION_URL       = "/application"
	APIV1                 = "/apiv1"
	ADMIN_ACCOUNTS        = APIV1 + "/admin/accounts"
	ADMIN_ACCOUNTS_CREATE = ADMIN_ACCOUNTS + "/create"
	ADMIN_GROUPS          = APIV1 + "/admin/groups"
	ADMIN_GROUPS_CREATE   = ADMIN_GROUPS + "/create"
	ADMIN_GROUPS_EDIT     = ADMIN_GROUPS + "/edit/{id}"
)
