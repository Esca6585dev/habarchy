package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Esca6585dev/habarchy/backend/api"
)

// swaggerUI is a minimal Swagger UI page loading the embedded spec.
const swaggerUI = `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Habarchy API</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css"></head>
<body><div id="ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
<script>window.ui = SwaggerUIBundle({url: "/api/docs/openapi.yaml", dom_id: "#ui", persistAuthorization: true, displayRequestDuration: true});</script>
</body></html>`

// RegisterDocs mounts GET /api/docs (Swagger UI) and /api/docs/openapi.yaml.
func RegisterDocs(app fiber.Router) {
	app.Get("/api/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString(swaggerUI)
	})
	app.Get("/api/docs/openapi.yaml", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
		return c.Send(api.OpenAPI)
	})
}
