package httpx

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// Exercise the public service constructor and router rather than a standalone
// Fiber app so dependency updates also check the forwarded configuration.
func TestServiceRoutingAndStatusCompatibility(t *testing.T) {
	config := GetDefaultConfig("127.0.0.1:0", "compatibility-test")
	config.CaseSensitive = true
	config.StrictRouting = true
	config.ErrorHandler = func(c fiber.Ctx, err error) error {
		status := http.StatusInternalServerError
		var httpErr *fiber.Error
		if errors.As(err, &httpErr) {
			status = httpErr.Code
		}
		return c.Status(status).SendString(err.Error())
	}
	service := NewService(config)
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	app := service.GetEngine()
	group := service.NewRouterGroup("/api")
	group.Get("/items/:id", func(c fiber.Ctx) error {
		return c.Status(http.StatusAccepted).SendString(c.Params("id") + ":" + c.Query("value"))
	})
	group.Get("/failure", func(fiber.Ctx) error {
		return fiber.NewError(http.StatusUnprocessableEntity, "invalid input")
	})
	for _, test := range []struct {
		name, method, path, body string
		status                   int
	}{
		{"parameter and query", http.MethodGet, "/api/items/one?value=two", "one:two", http.StatusAccepted},
		{"case sensitive", http.MethodGet, "/API/items/one", "", http.StatusNotFound},
		{"strict trailing slash", http.MethodGet, "/api/items/one/", "", http.StatusNotFound},
		{"method mismatch", http.MethodPost, "/api/items/one", "", http.StatusMethodNotAllowed},
		{"custom error status", http.MethodGet, "/api/failure", "invalid input", http.StatusUnprocessableEntity},
	} {
		t.Run(test.name, func(t *testing.T) {
			response, err := app.Test(httptest.NewRequest(test.method, test.path, http.NoBody))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil || response.StatusCode != test.status {
				t.Fatalf("status=%d body=%q err=%v", response.StatusCode, body, err)
			}
			if test.body != "" && string(body) != test.body {
				t.Fatalf("body=%q, want %q", body, test.body)
			}
		})
	}
}
