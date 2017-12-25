package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"github.com/aaronaaeng/chat.connor.fun/views"
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db/user"
	_"github.com/lib/pq"
	"github.com/labstack/echo/middleware"
)


func createApiRoutes(e *echo.Echo) {

}

func addMiddlewares(e *echo.Echo, c config.Config) {
	if !c.Debug {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}
}

func initDatabaseRepositories(c config.Config) {
	database, err := sqlx.Open("postgres", c.DatabaseURL)
	if err != nil {
		panic(err)
	}
	_, err = user.Init(database)
	if err != nil {
		panic(err)
	}
}

func main() {
	configData := config.New(true)

	initDatabaseRepositories(configData)

	e := echo.New()

	addMiddlewares(e, configData)

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}
	e.Renderer = t
	e.GET("/", views.Index)


	createApiRoutes(e)
	e.Logger.Fatal(e.Start(":4000"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}