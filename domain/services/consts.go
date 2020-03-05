package services

const (
	EMPTY = ""

	//group
	GORM_AUTOUPDATE            = "gorm:association_autoupdate"
	GORM_AUTOCREATE            = "gorm:association_autocreate"
	GROUP_ALREADY_EXIST        = "Group with this name already exist"
	ADMIN_REMOVE_PERMISSION    = "You can't remove Admin group"
	CUSTOMER_REMOVE_PERMISSION = "You can't remove Customer group"
	ADMIN_UPDATE_PERMISSION    = "You can't update Admin group"
	CUSTOMER_UPDATE_PERMISSION = "You can't update Customer group"
	GROUP_HAS_USERS            = "Group has users can't be removed"
	RECORD_NOT_FOUND           = "record not found"
	ADMIN                      = "Admin"
	CUSTOMER                   = "Customer"

	// user
	USER_ALREADY_EXIST     = "User with this mobile number already exist"
	USER_NOT_FOUND_WITH_ID = "User with this id not found"
)
