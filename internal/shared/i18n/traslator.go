package i18n

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/i18n/messages"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func T(locale, key string) string {
	localizer := goi18n.NewLocalizer(bundle, locale, "en")

	translated, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return key
	}
	return translated
}

var bundle = newBundle()

func newBundle() *goi18n.Bundle {
	i18nBundle := goi18n.NewBundle(language.English)

	i18nBundle.AddMessages(language.English, messages.GetENMessages()...)
	i18nBundle.AddMessages(language.BrazilianPortuguese, messages.GetPTBRMessages()...)

	return i18nBundle
}
