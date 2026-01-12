package handler

import (
	"net/http/httptest"
	"sync"
	"testing"

	"gopher-mock/model"

	"github.com/gofiber/fiber/v2"
)

func TestMockHandler_Concurrency(t *testing.T) {
	app := fiber.New()
	h := &MockHandler{
		Configs: []model.MockConfig{
			{
				Method: "GET",
				Path:   "/test",
				Responses: []model.ConditionalResponse{
					{
						Response: model.Response{
							StatusCode: 200,
							Body:       map[string]interface{}{"status": "ok"},
						},
					},
				},
			},
		},
	}

	app.All("/*", h.Dynamic)

	var wg sync.WaitGroup
	numRequests := 100

	// Concurrent reads
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/test", nil)
			_, _ = app.Test(req)
		}()
	}

	// Concurrent write
	wg.Add(1)
	go func() {
		defer wg.Done()
		newCfgs := []model.MockConfig{
			{
				Method: "GET",
				Path:   "/test",
				Responses: []model.ConditionalResponse{
					{
						Response: model.Response{
							StatusCode: 200,
							Body:       map[string]interface{}{"status": "updated"},
						},
					},
				},
			},
		}
		h.mu.Lock()
		h.Configs = newCfgs
		h.mu.Unlock()
	}()

	wg.Wait()
}
