package middlewares

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

func mapValidationTagToCode(tag string) string {
	switch tag {
	case "required":
		return shared_errors.CodeRequired
	case "email":
		return shared_errors.CodeInvalidEmail
	case "min":
		return shared_errors.CodeMinLength
	case "max":
		return shared_errors.CodeMaxLength
	default:
		return shared_errors.CodeInvalidFormat
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var errorDetails []shared_errors.ErrorDetail

		if errors.Is(err, io.EOF) {
			errorDetails = []shared_errors.ErrorDetail{
				{
					Code:  shared_errors.CodeMissingBody,
					Field: "",
				},
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errorDetails,
			})
			return
		}

		errMsg := err.Error()
		if strings.Contains(errMsg, "json") ||
			strings.Contains(errMsg, "bind") ||
			strings.Contains(errMsg, "unmarshal") ||
			strings.Contains(errMsg, "invalid character") {
			errorDetails = []shared_errors.ErrorDetail{
				{
					Code:  shared_errors.CodeInvalidJSON,
					Field: "",
				},
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errorDetails,
			})
			return
		}

		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			errorDetails = make([]shared_errors.ErrorDetail, 0, len(validationErrs))
			for _, fieldErr := range validationErrs {
				field := strings.ToLower(fieldErr.Field())
				tag := fieldErr.Tag()
				code := mapValidationTagToCode(tag)

				errorDetails = append(errorDetails, shared_errors.ErrorDetail{
					Code:  code,
					Field: field,
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errorDetails,
			})
			return
		}

		if appErr, ok := err.(*shared_errors.AppError); ok {
			errorDetail := appErr.ToErrorDetail()
			errorDetails = []shared_errors.ErrorDetail{errorDetail}

			c.JSON(appErr.Code, gin.H{
				"errors": errorDetails,
			})
			return
		}

		errorDetails = []shared_errors.ErrorDetail{
			{
				Code:  shared_errors.CodeInternalError,
				Field: "",
			},
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": errorDetails,
		})
	}
}