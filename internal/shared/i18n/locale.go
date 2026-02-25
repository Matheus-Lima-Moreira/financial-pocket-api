package i18n

import "strings"

func ResolveLocale(acceptLanguage string) string {
	lang := strings.ToLower(acceptLanguage)
	switch {
	case strings.HasPrefix(lang, "pt-br"), strings.HasPrefix(lang, "pt"):
		return "pt-BR"
	case strings.HasPrefix(lang, "en"):
		return "en"
	default:
		return "en"
	}
}
