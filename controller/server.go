package controller

import (
	"database/sql"
	"fmt"
	"log"

	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/dl-watson/pg-go/util"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"

	_ "github.com/lib/pq"
)

func setupRoutes(app *fiber.App, handler *Handler) {
	villagers := app.Group("/api/v1")

	setupCRUD(villagers, handler)
}

func setupViewRoutes(app *fiber.App, views *Views) {
	grp := app.Group("/views")
	setupViews(grp, views)
}

func setupCRUD(grp fiber.Router, handler *Handler) {
	routes := grp.Group("/villagers")
	routes.Get("/seed", handler.seedDB)
	routes.Get("/", handler.getVillagers)
	routes.Get("/:name", handler.getVillager)
	routes.Post("/", handler.createVillager)
}

func setupViews(grp fiber.Router, views *Views) {
	routes := grp.Group("/")
	routes.Get("/", views.villagersView)
	routes.Get("/:name", views.villagerView)
}

func SetupServer() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)

	engine := html.New("./views", ".html")
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})
	port := config.PORT

	handler := newHandler(store)
	views := newViews(store)

	setupRoutes(app, handler)
	setupViewRoutes(app, views)

	fmt.Println("Starting server on port 7890...")
	app.Listen(port)
}
