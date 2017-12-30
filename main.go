package main

import (
	"github.com/labstack/echo"
	"html/template"
	"io"
	"github.com/jmoiron/sqlx"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"github.com/aaronaaeng/chat.connor.fun/db/users"
	_"github.com/lib/pq"
	"github.com/labstack/echo/middleware"
	"github.com/aaronaaeng/chat.connor.fun/controllers"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"fmt"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
	_"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"io/ioutil"
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"strings"
	"github.com/aaronaaeng/chat.connor.fun/controllers/chat"
	"github.com/aaronaaeng/chat.connor.fun/context"
)


func createApiRoutes(e *echo.Echo) {
	e.POST("/api/v1/users", controllers.CreateUser)
	e.POST("/api/v1/login", controllers.LoginUser)

}

func addMiddlewares(e *echo.Echo) {
	if !config.Debug {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}
	//this must be added first
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.AuthorizedContextImpl{Context: c}
			return h(cc)
		}
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(jwtmiddleware.JwtAuth(func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/web") ||
				strings.HasPrefix(c.Path(), "/favicon.ico") || //skip static assets
				strings.HasSuffix(c.Path(), "ws")
	}, jwtmiddleware.JWTBearerTokenExtractor))

	e.Use(jwtmiddleware.JwtAuth(func(c echo.Context) bool { //websocket auth
		return !strings.HasSuffix(c.Path(), "ws")
	}, jwtmiddleware.JWTProtocolHeaderExtractor))
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
	e.Static("/web", "frontend")
	e.Debug = config.Debug

	roleJsonData, err := ioutil.ReadFile("assets/roles.json")
	if err != nil {
		e.Logger.Fatal(err)
	}
	if model.InitRoleMap(roleJsonData) != nil {
		e.Logger.Fatal(fmt.Errorf("error creating roles data"))
	}

	hubMap := chat.NewHubMap()

	initDatabaseRepositories()


	addMiddlewares(e)

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/public/*.html")),
	}
	e.Renderer = t
	e.GET("/", controllers.Index)
	e.GET("/wstest", controllers.WSTestView)

	e.GET("/api/v1/rooms/:room/messages/ws", func(c echo.Context) error {
		return chat.HandleWebsocket(hubMap, c)
	})

	createApiRoutes(e)
	e.Logger.Fatal(e.Start(":4000"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}