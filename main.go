package localewrapper

import (
	"encoding/json"
	"github.com/mindwingx/abstraction"
	"github.com/mindwingx/go-helper"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

type locale struct {
	bundle *i18n.Bundle
	Lang   string // Lang tag, like "en-US" and "fa-IR"
}

// NewLocale default language is set to English
func NewLocale(registry abstraction.Registry) abstraction.Locale {
	lang := new(locale)
	err := registry.Parse(&lang)
	if err != nil {
		helper.CustomPanic("[locale]registry init error:", err)
	}

	lang.bundle = i18n.NewBundle(language.English)

	return lang
}

// InitLocale language files formats could be contained JSON,YML, etc
func (l *locale) InitLocale(format string, locales []string) {
	l.bundle.RegisterUnmarshalFunc(format, json.Unmarshal)

	for _, localeFilePath := range locales {
		l.bundle.MustLoadMessageFile(localeFilePath)
	}
}

// Get returns message based on selected translation file
func (l *locale) Get(key string) string {
	localizer := i18n.NewLocalizer(l.bundle, l.Lang)

	localizedMessage, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key, // other fields are available with the i18n.Message struct
		},
	})

	return localizedMessage
}

// Plural returns message by params and multiple condition
func (l *locale) Plural(key string, params map[string]string) string {
	localizer := i18n.NewLocalizer(l.bundle, l.Lang)
	data := make(map[string]string)

	for localizerKey, localizerValue := range params {
		data[localizerKey] = localizerValue
	}

	formattedLocalizer := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: key,
		},
		TemplateData: data,
	})

	return formattedLocalizer
}

func (l *locale) FormatNumber(number int64) string {
	lang, _ := language.Parse(l.Lang) // Parse the language tag
	p := message.NewPrinter(lang)     // Create a new message printer with the specified language
	// Formatted number with grouping separators according to the user's preferred language
	return p.Sprintf("%d", number)
}

func (l *locale) FormatDate(date time.Time) string {
	lang, _ := language.Parse(l.Lang)
	p := message.NewPrinter(lang)
	// Format the date
	return p.Sprintf(
		"%s, %s %d, %d",
		date.Weekday(), date.Month(), date.Day(), date.Year(),
	)
}

func (l *locale) FormatCurrency(value float64, cur currency.Unit) string {
	lang, _ := language.Parse(l.Lang)
	p := message.NewPrinter(lang)
	// Format the currency value
	return p.Sprintf("%s %.2f", currency.Symbol(cur), value)
}
