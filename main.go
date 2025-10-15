package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"

	"golang-mock/handler"
	"golang-mock/template"
)

func main() {
	engine := html.New("./templates", ".html")
	engine.AddFunc("toJsonPretty", template.ToJSONPretty)
	engine.AddFunc("httpMethods", template.HTTPMethods)
	engine.AddFunc("len", template.ConfigLen)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	h := handler.NewMockHandler("configs.json")

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Get("/", h.Index)
	app.Post("/save", h.Save)
	app.Post("/delete-config/:index", h.Delete)
	app.All("/*", h.RequestResponseLogger(), h.Dynamic)

	log.Fatal(app.Listen(":3000"))
}
