package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"github.com/aaronaaeng/chat.connor.fun/views"
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db/users"
	_"github.com/lib/pq"
	"github.com/labstack/echo/middleware"
	_"github.com/mattn/go-sqlite3"
	"github.com/aaronaaeng/chat.connor.fun/controllers"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"fmt"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
)


func createApiRoutes(e *echo.Echo) {
	e.POST("/api/v1/users", controllers.CreateUser)
	e.POST("/api/v1/login", controllers.LoginUser)
}

func addMiddlewares(e *echo.Echo, c config.Config) {
	if !c.Debug {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}
}

func initDatabaseRepositories(c config.Config) {
	database, err := sqlx.Open("postgres", "postgresql://localhost:5432?sslmode=disable")
	if err != nil {
		panic(err)
	}
	_, err = users.Init(database)
	if err != nil {
		panic(err)
	}

	_, err = roles.Init(database)
	if err != nil {
		panic(err)
	}
}

func main() {
	configData := config.New(true)

	if model.InitRoleMap("assets/roles.json") != nil {
		fmt.Printf("Error creating User Roles mapping!")
	}

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