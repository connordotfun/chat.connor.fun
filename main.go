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
	"github.com/aaronaaeng/chat.connor.fun/db"
)


func createApiRoutes(api *echo.Group, hubMap *chat.HubMap, userRepository db.UserRepository,
	rolesRepository db.RolesRepository, roomsRepository db.RoomsRepository) {

	api.POST("/users", controllers.CreateUser(userRepository, rolesRepository)).Name = "create-user"
	api.GET("/users/:id", controllers.GetUser(userRepository)).Name = "get-user"

	api.POST("/login", controllers.LoginUser(userRepository)).Name = "login-user"

	api.GET("/rooms/*/ws", chat.HandleWebsocket(hubMap)).Name = "join-room"
}

func addMiddlewares(e *echo.Echo, rolesRepository db.RolesRepository) {
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
	}, jwtmiddleware.JWTBearerTokenExtractor, rolesRepository))

	e.Use(jwtmiddleware.JwtAuth(func(c echo.Context) bool { //websocket auth
		return !strings.HasSuffix(c.Path(), "ws")
	}, jwtmiddleware.JWTProtocolHeaderExtractor, rolesRepository))
}

func initDatabaseRepositories() (db.UserRepository, db.RolesRepository, db.RoomsRepository, db.MessagesRepository){
	database, err := sqlx.Open("postgres", config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	userRepository, err := users.New(database)
	if err != nil {
		panic(err)
	}

	rolesRepository, err := roles.New(database)
	if err != nil {
		panic(err)
	}

	return userRepository, rolesRepository, nil, nil
}

func main() {
	e := echo.New()
	e.Debug = config.Debug

	e.Static("/web", "frontend")
	e.GET("/", controllers.Index)
	e.GET("/wstest", controllers.WSTestView)
	v1ApiGroup := e.Group("/api/v1")

	roleJsonData, err := ioutil.ReadFile("assets/roles.json")
	if err != nil {
		e.Logger.Fatal(err)
	}
	if model.InitRoleMap(roleJsonData) != nil {
		e.Logger.Fatal(fmt.Errorf("error creating roles data"))
	}

	hubMap := chat.NewHubMap()
	usersRepository, rolesRepository, _, _ := initDatabaseRepositories()
	addMiddlewares(e, rolesRepository)
	createApiRoutes(v1ApiGroup, hubMap, usersRepository, rolesRepository, nil)

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/public/*.html")),
	}
	e.Renderer = t


	e.Logger.Fatal(e.Start(":4000"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}