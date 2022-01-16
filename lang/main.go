package lang

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strings"
)

type Service struct {
	bundle    *i18n.Bundle
	ctx       *gin.Context
	localizer *i18n.Localizer
}

func New(ctx *gin.Context, bundle *i18n.Bundle) Service {
	localizer := i18n.NewLocalizer(bundle, domainLanguage(ctx))
	return Service{
		bundle:    bundle,
		ctx:       ctx,
		localizer: localizer,
	}
}

func (s *Service) Trans(str string) string {
	// TODO, modiy this to handle
	for _, m := range translationMessages {
		if m.ID == str {
			localizedString, _ := s.localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &m,
			})
			return localizedString
		} else if m.Other == str {
			localizedString, _ := s.localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &m,
			})
			return localizedString
		}
	}
	return str
}

func domainLanguage(c *gin.Context) string {
	if strings.Contains(c.Request.Host, "turnering.io") {
		return "se"
	} else if strings.Contains(c.Request.Host, "turnering.local") {
		return "se"
	}
	return "en"
}
