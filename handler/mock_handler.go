package handler

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"golang-mock/config"
	"golang-mock/model"
	"golang-mock/service"
)

// MockHandler ...
type MockHandler struct {
	Configs []model.MockConfig
	Path    string
}

// NewMockHandler ...
func NewMockHandler(path string) *MockHandler {
	cfgs, err := config.LoadConfigs(path)
	if err != nil {
		log.Println("Load config error:", err)
	}
	return &MockHandler{
		Configs: cfgs,
		Path:    path,
	}
}

// Index ...
func (h *MockHandler) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Configs": h.Configs,
	})
}

// Save ...
func (h *MockHandler) Save(c *fiber.Ctx) error {
	var newCfgs []model.MockConfig
	if err := c.BodyParser(&newCfgs); err != nil {
		return c.Status(400).SendString("Invalid config format")
	}
	if err := config.SaveConfigs(h.Path, newCfgs); err != nil {
		return c.Status(500).SendString("Failed to save configs")
	}
	h.Configs = newCfgs
	return c.SendString("Config saved")
}

// Delete ...
func (h *MockHandler) Delete(c *fiber.Ctx) error {
	index, err := strconv.Atoi(c.Params("index"))
	if err != nil || index < 0 || index >= len(h.Configs) {
		return c.Status(400).SendString("Invalid index")
	}
	h.Configs = append(h.Configs[:index], h.Configs[index+1:]...)
	if err := config.SaveConfigs(h.Path, h.Configs); err != nil {
		return c.Status(500).SendString("Failed to save config")
	}
	return c.Redirect("/")
}

// Dynamic ...
func (h *MockHandler) Dynamic(c *fiber.Ctx) error {
	method := c.Method()
	path := c.Path()

	for _, cfg := range h.Configs {
		if cfg.Method == method && cfg.Path == path {
			// Validate headers
			headerMap := map[string]string{}
			for k := range cfg.RequestHeaders {
				val := c.Get(k)
				if val == "" {
					return c.Status(400).JSON(fiber.Map{"error": "missing header: " + k})
				}
				headerMap[k] = val
			}

			bodyMap := map[string]interface{}{}
			if (method == "POST" || method == "PUT") && cfg.RequestBody != nil {
				if err := c.BodyParser(&bodyMap); err != nil {
					return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
				}
				for k := range cfg.RequestBody {
					if _, ok := bodyMap[k]; !ok {
						return c.Status(400).JSON(fiber.Map{"error": "missing body field: " + k})
					}
				}
			}
			queryMap := make(map[string]string)
			c.Request().URI().QueryArgs().VisitAll(func(k, v []byte) {
				queryMap[string(k)] = string(v)
			})

			rendered := service.RenderTemplateRecursive(cfg.ResponseBody, service.MapToStringMap(bodyMap), headerMap, queryMap)

			for k, v := range cfg.ResponseHeaders {
				c.Set(k, v)
			}
			if cfg.Timeout > 0 {
				time.Sleep(time.Duration(cfg.Timeout) * time.Millisecond)
			}

			return c.Status(cfg.StatusCode).JSON(rendered)
		}
	}

	return c.Status(404).SendString("Mock not found")
}
