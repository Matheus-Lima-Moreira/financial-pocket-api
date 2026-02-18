package middlewares

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// Trata erro quando não há body (EOF)
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"corpo da requisição é obrigatório"},
			})
			return
		}

		// Trata erros de binding do Gin (JSON inválido, etc)
		errMsg := err.Error()
		if strings.Contains(errMsg, "json") ||
			strings.Contains(errMsg, "bind") ||
			strings.Contains(errMsg, "unmarshal") ||
			strings.Contains(errMsg, "invalid character") {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"formato JSON inválido"},
			})
			return
		}

		// Trata erros de validação do validator
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			var messages []string
			for _, fieldErr := range validationErrs {
				field := strings.ToLower(fieldErr.Field())
				tag := fieldErr.Tag()

				var message string
				switch tag {
				case "required":
					message = field + " é obrigatório"
				case "email":
					message = field + " deve ser um email válido"
				case "min":
					message = fmt.Sprintf("%s deve ter no mínimo %s caracteres", field, fieldErr.Param())
				case "max":
					message = fmt.Sprintf("%s deve ter no máximo %s caracteres", field, fieldErr.Param())
				default:
					message = fmt.Sprintf("%s é inválido", field)
				}
				messages = append(messages, message)
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": messages,
			})
			return
		}

		// Trata AppError customizado
		if appErr, ok := err.(*shared_errors.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"errors": []string{appErr.Message},
			})
			return
		}

		// Log do erro não tratado para debug
		fmt.Println("Erro não tratado:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": []string{"internal server error"},
		})
	}
}