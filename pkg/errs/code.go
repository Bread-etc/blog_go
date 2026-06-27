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

	// Aggregate 模块
	CodeAggregateTopPostsLimitInvalid = "AGGREGATE_TOP_POSTS_LIMIT_INVALID"

	// User 模块
	CodeInvalidCredentials     = "INVALID_CREDENTIALS"
	CodeInvalidPasswordEncrypt = "INVALID_PASSWORD_ENCRYPTION"
	CodeIncorrectOldPassword   = "INCORRECT_OLD_PASSWORD"
	CodeUsernameTooLong        = "USERNAME_TOO_LONG"
	CodePasswordTooShort       = "PASSWORD_TOO_SHORT"
	CodeUsernameRequired       = "USERNAME_REQUIRED"
	CodePasswordRequired       = "PASSWORD_REQUIRED"
	CodeOldPasswordRequired    = "OLD_PASSWORD_REQUIRED"
	CodeNewPasswordRequired    = "NEW_PASSWORD_REQUIRED"

	// Category 模块
	CodeCategoryNotFound     = "CATEGORY_NOT_FOUND"
	CodeCategoryNameExists   = "CATEGORY_NAME_EXISTS"
	CodeCategorySlugExists   = "CATEGORY_SLUG_EXISTS"
	CodeCategoryInUse        = "CATEGORY_IN_USE"
	CodeCategoryNameRequired = "CATEGORY_NAME_REQUIRED"
	CodeCategorySlugRequired = "CATEGORY_SLUG_REQUIRED"
	CodeCategoryNameTooLong  = "CATEGORY_NAME_TOO_LONG"
	CodeCategorySlugTooLong  = "CATEGORY_SLUG_TOO_LONG"

	// Tag 模块
	CodeTagNotFound     = "TAG_NOT_FOUND"
	CodeTagNameExists   = "TAG_NAME_EXISTS"
	CodeTagSlugExists   = "TAG_SLUG_EXISTS"
	CodeTagInUse        = "TAG_IN_USE"
	CodeTagNameRequired = "TAG_NAME_REQUIRED"
	CodeTagSlugRequired = "TAG_SLUG_REQUIRED"
	CodeTagNameTooLong  = "TAG_NAME_TOO_LONG"
	CodeTagSlugTooLong  = "TAG_SLUG_TOO_LONG"

	// Post 模块
	CodePostNotFound            = "POST_NOT_FOUND"
	CodePostSlugExists          = "POST_SLUG_EXISTS"
	CodePostCategoryNotFound    = "POST_CATEGORY_NOT_FOUND"
	CodePostTagNotFound         = "POST_TAG_NOT_FOUND"
	CodePostTitleRequired       = "POST_TITLE_REQUIRED"
	CodePostContentRequired     = "POST_CONTENT_REQUIRED"
	CodePostSlugRequired        = "POST_SLUG_REQUIRED"
	CodePostCategoryIDRequired  = "POST_CATEGORY_ID_REQUIRED"
	CodePostTagIDsRequired      = "POST_TAG_IDS_REQUIRED"
	CodePostIsPublishedRequired = "POST_IS_PUBLISHED_REQUIRED"
	CodePostTitleTooLong        = "POST_TITLE_TOO_LONG"
	CodePostSummaryTooLong      = "POST_SUMMARY_TOO_LONG"
	CodePostSlugTooLong         = "POST_SLUG_TOO_LONG"
	CodePostCoverTooLong        = "POST_COVER_TOO_LONG"
	CodePostKeywordTooLong      = "POST_KEYWORD_TOO_LONG"
	CodePostCoverURLInvalid     = "POST_COVER_URL_INVALID"
	CodePostCategoryIDInvalid   = "POST_CATEGORY_ID_INVALID"
	CodePostTagIDsInvalid       = "POST_TAG_IDS_INVALID"
	CodePostPageInvalid         = "POST_PAGE_INVALID"
	CodePostPageSizeInvalid     = "POST_PAGE_SIZE_INVALID"

	// Link 模块
	CodeLinkNotFound           = "LINK_NOT_FOUND"
	CodeLinkNameRequired       = "LINK_NAME_REQUIRED"
	CodeLinkURLRequired        = "LINK_URL_REQUIRED"
	CodeLinkURLInvalid         = "LINK_URL_INVALID"
	CodeLinkSortInvalid        = "LINK_SORT_INVALID"
	CodeLinkNameTooLong        = "LINK_NAME_TOO_LONG"
	CodeLinkURLTooLong         = "LINK_URL_TOO_LONG"
	CodeLinkDescriptionTooLong = "LINK_DESCRIPTION_TOO_LONG"

	// SiteConfig 模块
	CodeConfigTitleRequired      = "CONFIG_TITLE_REQUIRED"
	CodeConfigTitleTooLong       = "CONFIG_TITLE_TOO_LONG"
	CodeConfigSubtitleTooLong    = "CONFIG_SUBTITLE_TOO_LONG"
	CodeConfigDescriptionTooLong = "CONFIG_DESCRIPTION_TOO_LONG"
	CodeConfigKeywordsTooLong    = "CONFIG_KEYWORDS_TOO_LONG"
	CodeConfigAuthorTooLong      = "CONFIG_AUTHOR_TOO_LONG"
	CodeConfigEmailTooLong       = "CONFIG_EMAIL_TOO_LONG"
	CodeConfigGithubURLTooLong   = "CONFIG_GITHUB_URL_TOO_LONG"
	CodeConfigEmailInvalid       = "CONFIG_EMAIL_INVALID"
	CodeConfigGithubURLInvalid   = "CONFIG_GITHUB_URL_INVALID"

	// Database 模块
	CodeDatabaseUnavailable = "DATABASE_UNAVAILABLE"
)
