package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/gofiber/fiber"
)

type Handler struct {
	Store *db.Store
}

func newHandler(store *db.Store) *Handler {
	return &Handler{
		Store: store,
	}
}

func (h *Handler) seedDB(ctx *fiber.Ctx) {
	resp, err := http.Get("https://ac-vill.herokuapp.com/villagers?perPage=391")
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	var villager []db.CreateVillagerParams
	err = json.Unmarshal(body, &villager)
	if err != nil {
		log.Fatal(err)
	}

	for _, elem := range villager {
		villager, err := h.Store.CreateVillager(ctx.Context(), elem)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
			return
		}

		if err := ctx.Status(fiber.StatusCreated).JSON(villager); err != nil {
			ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
			return
		}
	}

	ctx.JSON(villager)
}

func (h *Handler) getVillager(ctx *fiber.Ctx) {
	villager, err := h.Store.GetVillager(ctx.Context(), ctx.Params("name"))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.JSON(villager)
}

func (h *Handler) getVillagers(ctx *fiber.Ctx) {
	villager, err := h.Store.GetVillagers(ctx.Context(), 397)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.JSON(villager)
}

func (h *Handler) createVillager(ctx *fiber.Ctx) {
	req := new(db.CreateVillagerParams)
	err := ctx.BodyParser(req)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "json err",
		})
		return
	}

	villager, err := h.Store.CreateVillager(ctx.Context(), *req)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(villager); err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return
	}

	ctx.JSON(villager)
}
