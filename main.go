package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type MockConfig struct {
	Name            string                 `json:"name"`
	Method          string                 `json:"method"`
	Path            string                 `json:"path"`
	RequestHeaders  map[string]string      `json:"requestHeaders"`
	RequestBody     map[string]interface{} `json:"requestBody"`
	ResponseHeaders map[string]string      `json:"responseHeaders"`
	ResponseBody    map[string]interface{} `json:"responseBody"`
	StatusCode      int                    `json:"statusCode"`
}

var dynamicGroup fiber.Router

var configs []MockConfig

func loadConfigs() error {
	data, err := os.ReadFile("configs.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &configs)
}

func saveConfigs(newConfigs []MockConfig) error {
	data, err := json.MarshalIndent(newConfigs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("configs.json", data, 0644)
}

func toJsonPretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

func setupDynamicRoutes(group fiber.Router) {
	for _, cfg := range configs {
		route := cfg // prevent closure bug
		group.Add(strings.ToUpper(route.Method), route.Path, func(c *fiber.Ctx) error {
			// Header validation
			for key, val := range route.RequestHeaders {
				if c.Get(key) != val {
					return c.Status(fiber.StatusBadRequest).SendString("Invalid header")
				}
			}

			// Body validation
			if len(route.RequestBody) > 0 {
				var body map[string]interface{}
				if err := c.BodyParser(&body); err != nil {
					return c.Status(fiber.StatusBadRequest).SendString("Invalid body")
				}
			}

			// Set response headers
			for key, val := range route.ResponseHeaders {
				c.Set(key, val)
			}

			return c.Status(route.StatusCode).JSON(route.ResponseBody)
		})
	}
}

func saveConfigToFile(configs []MockConfig) error {
	b, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("configs.json", b, 0644)
}

func main() {
	engine := html.New("./templates", ".html")
	engine.AddFunc("toJsonPretty", toJsonPretty)
	engine.AddFunc("len", func(arr []MockConfig) int { return len(arr) })

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	if err := loadConfigs(); err != nil {
		log.Println("No config loaded:", err)
		configs = []MockConfig{}
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Configs": configs,
		})
	})

	app.Post("/save", func(c *fiber.Ctx) error {
		var newConfigs []MockConfig
		if err := c.BodyParser(&newConfigs); err != nil {
			return c.Status(400).SendString("Invalid config format")
		}

		if err := saveConfigs(newConfigs); err != nil {
			return c.Status(500).SendString("Failed to write config")
		}

		configs = newConfigs
		dynamicGroup = app.Group("/_mock") // re-init group to remove old handlers
		setupDynamicRoutes(dynamicGroup)
		return c.SendString("Config saved successfully!")
	})

	app.Post("/delete-config/:index", func(c *fiber.Ctx) error {
		indexStr := c.Params("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return c.Status(400).SendString("Invalid index")
		}

		// configs, err := loadConfigs()
		// if err != nil {
		// 	return c.Status(500).SendString("Gagal membaca konfigurasi")
		// }

		if index < 0 || index >= len(configs) {
			return c.Status(400).SendString("Index out of range")
		}

		// ❌ Ini salah jika kamu tidak menyimpan array baru
		// configs = append(configs[:index], configs[index+1:]...)

		// ✅ Ini benar
		newConfigs := append(configs[:index], configs[index+1:]...)

		// Simpan ulang
		if err := saveConfigs(newConfigs); err != nil {
			return c.Status(500).SendString("Gagal menyimpan konfigurasi")
		}
		dynamicGroup = app.Group("/_mock") // re-init group to remove old handlers
		setupDynamicRoutes(dynamicGroup)

		return c.Redirect("/")
	})

	dynamicGroup = app.Group("/_mock")
	setupDynamicRoutes(dynamicGroup)

	log.Fatal(app.Listen(":3000"))
}
