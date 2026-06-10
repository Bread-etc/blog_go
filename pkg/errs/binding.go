package errs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FromBinding(err error) *Error {
	if err == nil {
		return nil
	}

	// 检测是否为 validator 校验器错误（默认错误 + formValidationError定义错误）
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		if len(validationErrors) == 0 {
			return New(http.StatusBadRequest, CodeInvalidParams, "invalid request parameters")
		}
		return fromValidationError(validationErrors[0])
	}

	// 检测是否为 syntax 语法错误
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return New(http.StatusBadRequest, CodeInvalidParams, "invalid request body")
	}

	// 检测是否为 type 类型错误
	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		return New(http.StatusBadRequest, CodeInvalidParams, "invalid request parameter type")
	}

	// 检测特定错误值是否为空（未填写）
	if errors.Is(err, io.EOF) {
		return New(http.StatusBadRequest, CodeInvalidParams, "request body is required")
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid request parameters")
}

func fromValidationError(fe validator.FieldError) *Error {
	namespace := fe.StructNamespace()
	field := fe.StructField()
	tag := fe.Tag()

	switch {
	case strings.Contains(namespace, "Category"):
		return fromCategoryValidationError(field, tag)
	case strings.Contains(namespace, "Tag"):
		return fromTagValidationError(field, tag)
	case strings.Contains(namespace, "Post"):
		return fromPostValidationError(field, tag)
	case strings.Contains(namespace, "Link"):
		return fromLinkValidationError(field, tag)
	case strings.Contains(namespace, "Config"):
		return fromConfigValidationError(field, tag)
	case strings.Contains(namespace, "LoginReq"), strings.Contains(namespace, "ChangePassword"):
		return fromUserValidationError(field, tag)
	default:
		return New(http.StatusBadRequest, CodeInvalidParams, "invalid request parameters")
	}
}

func fromCategoryValidationError(field string, tag string) *Error {
	switch field {
	case "Name":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeCategoryNameRequired, "category name is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeCategoryNameTooLong, "category name is too long")
		}
	case "Slug":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeCategorySlugRequired, "category slug is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeCategorySlugTooLong, "category slug is too long")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid category parameter")
}

func fromTagValidationError(field string, tag string) *Error {
	switch field {
	case "Name":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeTagNameRequired, "tag name is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeTagNameTooLong, "tag name is too long")
		}
	case "Slug":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeTagSlugRequired, "tag slug is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeTagSlugTooLong, "tag slug is too long")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid tag parameters")
}

func fromPostValidationError(field string, tag string) *Error {
	switch field {
	case "Title":
		if tag == "required" {
			return New(http.StatusBadRequest, CodePostTitleRequired, "post title is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodePostTitleTooLong, "post title is too long")
		}
	case "Content":
		if tag == "required" {
			return New(http.StatusBadRequest, CodePostContentRequired, "post content is required")
		}
	case "Slug":
		if tag == "required" {
			return New(http.StatusBadRequest, CodePostSlugRequired, "post slug is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodePostSlugTooLong, "post slug is too long")
		}
	case "Summary":
		if tag == "max" {
			return New(http.StatusBadRequest, CodePostSummaryTooLong, "post summary is too long")
		}
	case "Cover":
		if tag == "url" {
			return New(http.StatusBadRequest, CodePostCoverURLInvalid, "post cover url is invalid")
		}
	case "CategoryID":
		if tag == "required" {
			return New(http.StatusBadRequest, CodePostCategoryIDRequired, "post category id is required")
		}
	case "TagIDs":
		if tag == "required" || tag == "min" {
			return New(http.StatusBadRequest, CodePostTagIDsRequired, "post tag ids are required")
		}
	case "Page":
		if tag == "min" {
			return New(http.StatusBadRequest, CodePostPageInvalid, "post page is invalid")
		}
	case "PageSize":
		if tag == "min" || tag == "max" {
			return New(http.StatusBadRequest, CodePostPageSizeInvalid, "post page size is invalid")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid post parameters")
}

func fromLinkValidationError(field string, tag string) *Error {
	switch field {
	case "Name":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeLinkNameRequired, "link name is required")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeLinkNameTooLong, "link name is too long")
		}
	case "URL":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeLinkURLRequired, "link url is required")
		}
		if tag == "url" {
			return New(http.StatusBadRequest, CodeLinkURLInvalid, "link url is invalid")
		}
		if tag == "max" {
			return New(http.StatusBadRequest, CodeLinkURLTooLong, "link url is too long")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid link parameters")
}

func fromConfigValidationError(field string, tag string) *Error {
	switch field {
	case "Title":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeConfigTitleRequired, "config title is required")
		}
	case "Email":
		if tag == "email" {
			return New(http.StatusBadRequest, CodeConfigEmailInvalid, "config email is invalid")
		}
	case "GithubURL":
		if tag == "url" {
			return New(http.StatusBadRequest, CodeConfigGithubURLInvalid, "config github url is invalid")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid config parameters")
}

func fromUserValidationError(field string, tag string) *Error {
	switch field {
	case "Username":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeUsernameRequired, "username is required")
		}
	case "Password":
		if tag == "required" {
			return New(http.StatusBadRequest, CodePasswordRequired, "password is required")
		}
	case "OldPassword":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeOldPasswordRequired, "old password is required")
		}
	case "NewPassword":
		if tag == "required" {
			return New(http.StatusBadRequest, CodeNewPasswordRequired, "new password is required")
		}
		if tag == "min" {
			return New(http.StatusBadRequest, CodePasswordTooShort, "password is too short")
		}
	}

	return New(http.StatusBadRequest, CodeInvalidParams, "invalid user parameters")
}
