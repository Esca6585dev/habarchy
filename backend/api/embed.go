// Package api embeds the OpenAPI document so the server can serve it at
// /api/docs together with Swagger UI.
package api

import _ "embed"

// OpenAPI is the OpenAPI 3.1 document (api/openapi.yaml).
//
//go:embed openapi.yaml
var OpenAPI []byte
