package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"github.com/aaronaaeng/chat.connor.fun/views"
)


func createApiRoutes(e *echo.Echo) {

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}
	e.Renderer = t
	e.GET("/", views.Index)


	createApiRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}