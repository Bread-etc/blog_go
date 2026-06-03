package errs

// 业务错误码字符串定义

const (
	CodeSuccess       = ""               // 200
	CodeInvalidParams = "INVALID_PARAMS" // 400
	CodeUnauthorized  = "UNAUTHORIZED"   // 401
	CodeForbidden     = "FORBIDDEN"      // 403
	CodeNotFound      = "NOT_FOUND"      // 404
	CodeConflict      = "CONFLICT"       // 409
	CodeInternalError = "INTERNAL_ERROR" // 500

	// User 模块
	CodeInvalidCredentials     = "INVALID_CREDENTIALS"
	CodeInvalidPasswordEncrypt = "INVALID_PASSWORD_ENCRYPTION"
	CodeIncorrectOldPassword   = "INCORRECT_OLD_PASSWORD"
	CodePasswordTooShort       = "PASSWORD_TOO_SHORT"

	// Category 模块
	CodeCategoryNotFound   = "CATEGORY_NOT_FOUND"
	CodeCategoryNameExists = "CATEGORY_NAME_EXISTS"
	CodeCategorySlugExists = "CATEGORY_SLUG_EXISTS"
	CodeCategoryInUse      = "CATEGORY_IN_USE"

	// Tag 模块
	CodeTagNotFound   = "TAG_NOT_FOUND"
	CodeTagNameExists = "TAG_NAME_EXISTS"
	CodeTagSlugExists = "TAG_SLUG_EXISTS"
	CodeTagInUse      = "TAG_IN_USE"

	// Post 模块
	CodePostNotFound         = "POST_NOT_FOUND"
	CodePostSlugExists       = "POST_SLUG_EXISTS"
	CodePostCategoryNotFound = "POST_CATEGORY_NOT_FOUND"
	CodePostTagNotFound      = "POST_TAG_NOT_FOUND"

	// Link 模块
	CodeLinkNotFound = "LINK_NOT_FOUND"

	// SiteConfig 模块
	CodeConfigInvalidEmail     = "CONFIG_INVALID_EMAIL"
	CodeConfigInvalidGithubURL = "CONFIG_INVALID_GITHUB_URL"

	// Database 模块
	CodeDatabaseUnavailable = "DATABASE_UNAVAILABLE"
)
