package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)


//GET('/')
func Index(c echo.Context) error {
	templateVars := map[string]interface{} {
		"publicUrl": "/web/public",
	}
	return c.Render(http.StatusOK, "index.html", templateVars)
}