package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
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

	app.Use(h.RequestResponseLogger())
	app.Get("/", h.Index)
	app.Post("/save", h.Save)
	app.Post("/delete-config/:index", h.Delete)
	app.All("/*", h.Dynamic)

	log.Fatal(app.Listen(":3000"))
}
