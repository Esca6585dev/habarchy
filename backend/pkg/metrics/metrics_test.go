package metrics

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func TestFiberMiddlewareLabels(t *testing.T) {
	app := fiber.New()
	app.Use(FiberMiddleware())
	app.Get("/metrics", adaptor.HTTPHandler(Handler()))
	app.Post("/api/v1/messages/:id", func(c *fiber.Ctx) error { return c.SendStatus(202) })
	app.Get("/boom", func(*fiber.Ctx) error { return fiber.ErrTeapot })

	for i := 0; i < 3; i++ {
		res, _ := app.Test(httptest.NewRequest("POST", "/api/v1/messages/abc-123", strings.NewReader("{}")))
		if res.StatusCode != 202 {
			t.Fatal(res.StatusCode)
		}
	}
	_, _ = app.Test(httptest.NewRequest("GET", "/boom", nil))
	_, _ = app.Test(httptest.NewRequest("GET", "/nowhere", nil))

	res, _ := app.Test(httptest.NewRequest("GET", "/metrics", nil))
	body, _ := io.ReadAll(res.Body)
	out := string(body)
	for _, want := range []string{
		`habarchy_http_requests_total{method="POST",route="/api/v1/messages/:id",status="2xx"} 3`,
		`habarchy_http_requests_total{method="GET",route="/boom",status="4xx"} 1`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, filter(out, "habarchy_http_requests_total"))
		}
	}
	if strings.Contains(out, "abc-123") {
		t.Error("raw path leaked into labels")
	}
}

func filter(s, prefix string) string {
	var b strings.Builder
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, prefix) {
			b.WriteString(l + "\n")
		}
	}
	return b.String()
}
