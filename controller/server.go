package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	db "github.com/dl-watson/pg-go/db/sqlc"
	"github.com/dl-watson/pg-go/util"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"

	_ "github.com/lib/pq"
)

type Handler struct {
	Store *db.Store
}

type Views struct {
	Store *db.Store
}

func NewHandler(store *db.Store) *Handler {
	return &Handler{
		Store: store,
	}
}
func NewViews(store *db.Store) *Views {
	return &Views{
		Store: store,
	}
}

func (h *Handler) SeedDB(ctx *fiber.Ctx) {
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

		fmt.Println(villager)
	}

	ctx.JSON(villager)
}

func (h *Handler) GetVillager(ctx *fiber.Ctx) {
	villager, err := h.Store.GetVillager(ctx.Context(), ctx.Params("name"))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.JSON(villager)
}

func (h *Handler) GetVillagers(ctx *fiber.Ctx) {
	villager, err := h.Store.GetVillagers(ctx.Context(), 397)
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
		ctx.Status(fiber.StatusInternalServerError)
		return
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(villager); err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return
	}

	ctx.JSON(villager)
}


func (v *Views) VillagersView(ctx *fiber.Ctx)  {
	villager, err := v.Store.GetVillagers(ctx.Context(), 397)
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.Render("list", fiber.Map{
		"Villager": villager,
	})

}

func (v *Views) VillagerView(ctx *fiber.Ctx)  {
	villager, err := v.Store.GetVillager(ctx.Context(), ctx.Params("name"))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	ctx.Render("detail", fiber.Map{
		"Name": villager.Name,
		"Image": villager.Image,
		"Species": villager.Species,
		"Personality": villager.Personality,
		"Birthday": villager.Birthday,
		"Quote": villager.Quote,
	})
}

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