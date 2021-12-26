package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type IndexData struct {
	PageData
	SiteNameFirst            string
	SiteNameLast             string
	EasyToCreateTournaments  string
	TournamentName           string
	Optional                 string
	Create                   string
	SubscribeToOurNewsletter string
	Email                    string
	Subscribe                string
	NoSpamNotice             string
}

func (controller Controller) Index(c *gin.Context) {
	localize := i18n.NewLocalizer(controller.bundle, domainLanguage(c))

	siteNameFirst, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "site_name_first",
			Other: "Tournify",
		},
	})

	siteNameLast, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "site_name_last",
			Other: ".io",
		},
	})

	easyToCreateTournaments, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "easy_to_create_tournaments",
			Other: "Easy to create tournaments",
		},
	})

	tournamentName, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "tournament_name",
			Other: "Tournament name",
		},
	})

	optional, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "optional",
			Other: "Optional",
		},
	})

	create, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "create",
			Other: "Create",
		},
	})

	subscribeToOurNewsletter, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "subscribe_to_our_newsletter",
			Other: "Subscribe to our newsletter",
		},
	})

	email, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "email",
			Other: "Email",
		},
	})

	subscribe, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "subscribe",
			Other: "Subscribe",
		},
	})

	noSpamNotice, _ := localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "no_spam_notice",
			Other: "If you want to receive updates from us just pop your email in the box. We think that spam is for jerks. And jerks we are not.",
		},
	})

	id := IndexData{
		PageData:                 controller.defaultPageData(c),
		SiteNameFirst:            siteNameFirst,
		SiteNameLast:             siteNameLast,
		EasyToCreateTournaments:  easyToCreateTournaments,
		TournamentName:           tournamentName,
		Optional:                 optional,
		Create:                   create,
		SubscribeToOurNewsletter: subscribeToOurNewsletter,
		Email:                    email,
		Subscribe:                subscribe,
		NoSpamNotice:             noSpamNotice,
	}
	c.HTML(http.StatusOK, "index.html", id)
}
