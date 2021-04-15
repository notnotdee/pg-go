package controller

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/dl-watson/pg-go/util"
	"github.com/gofiber/fiber"
)

type Handler struct {
	Store *db.Store
}

func NewHandler(store *db.Store) *Handler {
	return &Handler{
		Store: store,
	}
}

func (h *Handler) GetVillager(ctx *fiber.Ctx) {
	villager, err := h.Store.GetVillager(ctx.Context(), ctx.Params("name"))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.JSON(villager)
}

func (h *Handler) CreateVillager(ctx *fiber.Ctx) {
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
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(villager); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	ctx.JSON(villager)
}

func setupRoutes(app *fiber.App, handler *Handler) {
	path := app.Group("/api/v1")

	setupCRUD(path, handler)
}

func setupCRUD(grp fiber.Router, handler *Handler) {
	routes := grp.Group("/villagers")
	routes.Get("/:name", handler.GetVillager)
	routes.Post("/", handler.CreateVillager)
}

func SetupServer() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	
	app := fiber.New()
	port := config.PORT

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("hello world")
	})
	
	handler := NewHandler(store)

	setupRoutes(app, handler)

	fmt.Println("Starting server on port 7890...")
	app.Listen(port)
}