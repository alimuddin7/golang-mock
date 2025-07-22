package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// MockConfig ...
type MockConfig struct {
	Name            string                 `json:"name"`
	Method          string                 `json:"method"`
	Path            string                 `json:"path"`
	RequestHeaders  map[string]string      `json:"requestHeaders"`
	RequestBody     map[string]interface{} `json:"requestBody"`
	ResponseHeaders map[string]string      `json:"responseHeaders"`
	ResponseBody    map[string]interface{} `json:"responseBody"`
	StatusCode      int                    `json:"statusCode"`
	Timeout         int                    `json:"timeout"` // Timeout in miliseconds
}

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

func renderTemplateRecursive(data interface{}, bodyMap, headerMap, queryMap map[string]string) interface{} {
	switch v := data.(type) {
	case string:
		// replace {{key}} placeholders
		for k, val := range bodyMap {
			v = strings.ReplaceAll(v, "{{body."+k+"}}", val)
		}
		for k, val := range headerMap {
			v = strings.ReplaceAll(v, "{{header."+k+"}}", val)
		}
		for k, val := range queryMap {
			v = strings.ReplaceAll(v, "{{query."+k+"}}", val)
		}
		return v
	case map[string]interface{}:
		result := map[string]interface{}{}
		for key, val := range v {
			result[key] = renderTemplateRecursive(val, bodyMap, headerMap, queryMap)
		}
		return result
	case []interface{}:
		for i, val := range v {
			v[i] = renderTemplateRecursive(val, bodyMap, headerMap, queryMap)
		}
	}
	return data
}

func dynamicRouteHandler(c *fiber.Ctx) error {
	method := c.Method()
	path := c.Path()

	for _, cfg := range configs {
		if strings.EqualFold(cfg.Method, method) && cfg.Path == path {
			// Validate headers

			if cfg.StatusCode == 500 {
				return c.Status(500).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			}
			headerMap := map[string]string{}
			for key := range cfg.RequestHeaders {
				headerMap[key] = c.Get(key)
				if c.Get(key) == "" {
					return c.Status(400).JSON(fiber.Map{
						"error": "missing required header: " + key,
					})
				}
			}

			bodyMap := map[string]interface{}{}
			bodyStrMap := map[string]string{}
			if cfg.RequestBody != nil && (c.Method() == "POST" || c.Method() == "PUT") {
				if err := c.BodyParser(&bodyMap); err != nil {
					return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
				}
				for key := range cfg.RequestBody {
					if _, ok := bodyMap[key]; !ok {
						return c.Status(400).JSON(fiber.Map{
							"error": "missing required body field: " + key,
						})
					}
				}
				for k, v := range bodyMap {
					bodyStrMap[k] = fmt.Sprintf("%v", v)
				}
			}

			queryMap := make(map[string]string)
			c.Request().URI().QueryArgs().VisitAll(func(k, v []byte) {
				queryMap[string(k)] = string(v)
			})

			rendered := renderTemplateRecursive(cfg.ResponseBody, bodyStrMap, headerMap, queryMap)

			for k, v := range cfg.ResponseHeaders {
				c.Set(k, v)
			}
			if cfg.Timeout > 0 {
				time.Sleep(time.Duration(cfg.Timeout) * time.Millisecond)
			}

			return c.Status(cfg.StatusCode).JSON(rendered)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Mock not found")
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
			return c.Status(fiber.StatusBadRequest).SendString("Invalid config format")
		}

		if err := saveConfigs(newConfigs); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to write config")
		}

		configs = newConfigs
		return c.SendString("Config saved successfully!")
	})

	app.Post("/delete-config/:index", func(c *fiber.Ctx) error {
		indexStr := c.Params("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 0 || index >= len(configs) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid index")
		}

		newConfigs := append(configs[:index], configs[index+1:]...)

		if err := saveConfigs(newConfigs); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save configs")
		}

		configs = newConfigs
		return c.Redirect("/")
	})

	// Register universal catch-all route AFTER all others
	app.All("/*", dynamicRouteHandler)

	log.Fatal(app.Listen(":3000"))
}
