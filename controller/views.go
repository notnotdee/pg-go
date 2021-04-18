package controller

import (
	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/gofiber/fiber"
)

type Views struct {
	Store *db.Store
}

func newViews(store *db.Store) *Views {
	return &Views{
		Store: store,
	}
}

func (v *Views) villagersView(ctx *fiber.Ctx) {
	villager, err := v.Store.GetVillagers(ctx.Context(), 397)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.Render("list", fiber.Map{
		"Villager": villager,
	})
}

func (v *Views) villagerView(ctx *fiber.Ctx) {
	villager, err := v.Store.GetVillager(ctx.Context(), ctx.Params("name"))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.Render("detail", fiber.Map{
		"Name":        villager.Name,
		"Image":       villager.Image,
		"Species":     villager.Species,
		"Personality": villager.Personality,
		"Birthday":    villager.Birthday,
		"Quote":       villager.Quote,
	})
}
