package routes

import (
	"github.com/tournify/web/models"
	"github.com/tournify/web/util"
)

func (controller Controller) createUniqueTournamentSlug(slug string, iteration int) string {
	t := models.Tournament{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniqueTournamentSlug(slug, iteration+1)
}

func (controller Controller) createUniqueGameSlug(slug string, iteration int) string {
	t := models.Game{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniqueGameSlug(slug, iteration+1)
}

func (controller Controller) createUniqueGroupSlug(slug string, iteration int) string {
	t := models.Group{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniqueGroupSlug(slug, iteration+1)
}

func (controller Controller) createUniquePlayerSlug(slug string, iteration int) string {
	t := models.Player{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniquePlayerSlug(slug, iteration+1)
}

func (controller Controller) createUniquePostSlug(slug string, iteration int) string {
	t := models.Post{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniquePostSlug(slug, iteration+1)
}

func (controller Controller) createUniqueTeamSlug(slug string, iteration int) string {
	t := models.Team{
		Slug: slug,
	}
	controller.db.Where(t).First(&t)
	if t.ID == 0 {
		return slug
	}
	if iteration > 0 {
		slug = slug + util.RandomString(1)
	} else {
		slug = slug + "-" + util.RandomString(4)
	}
	return controller.createUniqueTeamSlug(slug, iteration+1)
}
