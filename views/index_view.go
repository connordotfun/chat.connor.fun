package views

import (
	"github.com/labstack/echo"
	"net/http"
)


//GET('/')
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}