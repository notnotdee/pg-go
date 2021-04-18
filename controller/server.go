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
	routes.Get("/seed", handler.SeedDB)
	routes.Get("/", handler.GetVillagers)
	routes.Get("/:name", handler.GetVillager)
	routes.Post("/", handler.CreateVillager)
}

func setupViews(grp fiber.Router, views *Views) {
	routes := grp.Group("/")
	routes.Get("/", views.VillagersView)
	routes.Get("/:name", views.VillagerView)
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

	handler := NewHandler(store)
	views := NewViews(store)

	setupRoutes(app, handler)
	setupViewRoutes(app, views)

	fmt.Println("Starting server on port 7890...")
	app.Listen(port)
}
