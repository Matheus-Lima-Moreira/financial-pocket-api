package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/i18n"
	"github.com/gin-gonic/gin"
)

type i18nResponseWriter struct {
	gin.ResponseWriter
	body        *bytes.Buffer
	statusCode  int
	wroteHeader bool
}

func newI18nResponseWriter(writer gin.ResponseWriter) *i18nResponseWriter {
	return &i18nResponseWriter{
		ResponseWriter: writer,
		body:           bytes.NewBuffer(nil),
		statusCode:     http.StatusOK,
	}
}

func (w *i18nResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = code
	w.wroteHeader = true
}

func (w *i18nResponseWriter) WriteHeaderNow() {
	w.wroteHeader = true
}

func (w *i18nResponseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(data)
}

func (w *i18nResponseWriter) WriteString(value string) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.WriteString(value)
}

func (w *i18nResponseWriter) Status() int {
	return w.statusCode
}

func (w *i18nResponseWriter) Size() int {
	return w.body.Len()
}

func (w *i18nResponseWriter) Written() bool {
	return w.wroteHeader || w.body.Len() > 0
}

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		originalWriter := c.Writer
		translatedWriter := newI18nResponseWriter(originalWriter)
		c.Writer = translatedWriter

		c.Next()

		body := translatedWriter.body.Bytes()
		body = translateReplyBody(c, body, translatedWriter.Header().Get("Content-Type"))

		c.Writer = originalWriter
		if translatedWriter.Written() {
			originalWriter.WriteHeader(translatedWriter.Status())
		}

		if len(body) > 0 {
			_, _ = originalWriter.Write(body)
		}
	}
}

func translateReplyBody(c *gin.Context, body []byte, contentType string) []byte {
	if len(body) == 0 {
		return body
	}

	if !strings.Contains(strings.ToLower(contentType), "application/json") {
		return body
	}

	var reply dtos.ReplyDTO
	if err := json.Unmarshal(body, &reply); err != nil {
		return body
	}

	if reply.Message == "" {
		return body
	}

	locale := i18n.ResolveLocale(c.GetHeader("Accept-Language"))
	reply.Message = i18n.T(locale, reply.Message)

	translatedBody, err := json.Marshal(reply)
	if err != nil {
		return body
	}

	return translatedBody
}
