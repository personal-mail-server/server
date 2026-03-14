package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const swaggerHTML = `<!doctype html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Personal Mail Server API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/docs/openapi.yaml',
        dom_id: '#swagger-ui'
      });
    </script>
  </body>
</html>
`

func DocsPage(c echo.Context) error {
	return c.HTML(http.StatusOK, swaggerHTML)
}
