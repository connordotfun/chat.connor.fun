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
	"github.com/aaronaaeng/chat.connor.fun/controllers"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"fmt"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"io/ioutil"
)


func createApiRoutes(e *echo.Echo) {
	e.POST("/api/v1/users", controllers.CreateUser)
	e.POST("/api/v1/login", controllers.LoginUser)
}

func addMiddlewares(e *echo.Echo) {
	if !config.Debug {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(jwtmiddleware.JwtAuth())
}

func initDatabaseRepositories() {
	database, err := sqlx.Open("postgres", config.DatabaseURL)
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
	e := echo.New()
	e.Debug = config.Debug

	roleJsonData, err := ioutil.ReadFile("assets/roles.json")
	if err != nil {
		e.Logger.Fatal(err)
	}
	if model.InitRoleMap(roleJsonData) != nil {
		e.Logger.Fatal(fmt.Errorf("error creating roles data"))
	}

	initDatabaseRepositories()


	addMiddlewares(e)

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