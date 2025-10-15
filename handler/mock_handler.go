package handler

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"golang-mock/config"
	"golang-mock/model"
	"golang-mock/service"
)

// MockHandler ...
type MockHandler struct {
	Configs []model.MockConfig
	Path    string
	Log     zerolog.Logger
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
		Log:     zerolog.New(os.Stdout).With().Timestamp().Logger(),
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
		if cfg.Method == method {
			if ok, params := pathMatch(cfg.Path, path); ok {
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

				rendered := service.RenderTemplateRecursive(cfg.ResponseBody, service.MapToStringMap(bodyMap), headerMap, queryMap, params)

				for k, v := range cfg.ResponseHeaders {
					c.Set(k, v)
				}
				if cfg.Timeout > 0 {
					time.Sleep(time.Duration(cfg.Timeout) * time.Millisecond)
				}

				return c.Status(cfg.StatusCode).JSON(rendered)
			}
		}
	}

	return c.Status(404).SendString("Mock not found")
}

func pathMatch(cfgPath, reqPath string) (bool, map[string]string) {
	cfgParts := strings.Split(strings.Trim(cfgPath, "/"), "/")
	reqParts := strings.Split(strings.Trim(reqPath, "/"), "/")

	if len(cfgParts) != len(reqParts) {
		return false, nil
	}

	params := map[string]string{}
	for i := range cfgParts {
		if strings.HasPrefix(cfgParts[i], ":") {
			paramName := strings.TrimPrefix(cfgParts[i], ":")
			params[paramName] = reqParts[i]
		} else if cfgParts[i] != reqParts[i] {
			return false, nil
		}
	}
	return true, params
}

// RequestResponseLogger ...
func (h *MockHandler) RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		reqHeaders := c.GetReqHeaders()
		reqBody := c.Body()

		err := c.Next()

		resBody := c.Response().Body()

		event := h.Log.Info().
			Str("METHOD", c.Method()).
			Str("URL", c.OriginalURL()).
			Any("HEAD", reqHeaders).
			Any("BODY", parseJSONOrString(reqBody)).
			Any("RES", parseJSONOrString(resBody)).
			Dur("DURATION", time.Since(start))

		if err != nil {
			event = event.Err(err)
		}

		event.Msg("API LOG")

		return err
	}
}

func parseJSONOrString(body []byte) any {
	if len(body) == 0 {
		return nil
	}

	var v any
	if err := json.Unmarshal(body, &v); err == nil {
		return v // berhasil parse JSON → return object/map
	}

	// kalau bukan JSON valid (misal text biasa)
	return string(body)
}
