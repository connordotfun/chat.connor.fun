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
	"io/ioutil"
	"github.com/aaronaaeng/chat.connor.fun/controllers/jwtmiddleware"
	"strings"
	"github.com/aaronaaeng/chat.connor.fun/controllers/chat"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/aaronaaeng/chat.connor.fun/db/rooms"
	"github.com/aaronaaeng/chat.connor.fun/db/messages"
	"github.com/aaronaaeng/chat.connor.fun/db/verifications"
)


var (
	userRepository db.UserRepository
	rolesRepository db.RolesRepository
	roomsRepository db.RoomsRepository
	messagesRepository db.MessagesRepository
	verificationsRepository db.VerificationCodeRepository
)

func createApiRoutes(api *echo.Group, hubMap *chat.HubMap) {

	api.POST("/users", controllers.CreateUser(userRepository, rolesRepository, verificationsRepository)).Name = "create-user"
	api.GET("/users/:id", controllers.GetUser(userRepository)).Name = "get-user"
	api.PUT("/users/:id", controllers.UpdateUser(userRepository))

	api.POST("/login", controllers.LoginUser(userRepository)).Name = "login-user"

	api.GET("/messages", controllers.GetMessages(messagesRepository)).Name = "get-messages"
	api.GET("/messages/:id", controllers.GetMessage(messagesRepository)).Name = "get-message"
	api.PUT("/messages/:id", controllers.UpdateMessage(messagesRepository)).Name = "update-message"

	api.GET("/rooms/nearby", controllers.GetNearbyRooms(roomsRepository)).Name = "get-nearby-rooms"
	api.GET("/rooms/:room/users", controllers.GetRoomMembers(hubMap)).Name = "get-room-members"
	api.GET("/rooms/:room", controllers.GetRoom(roomsRepository, hubMap)).Name = "get-room"

	api.PUT("/verifications/accountverification", controllers.VerifyUserAccount(verificationsRepository, userRepository, rolesRepository))

	api.GET("/rooms/:room/ws", chat.HandleWebsocket(hubMap, roomsRepository, messagesRepository)).Name = "join-room"
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
			strings.HasPrefix(c.Path(), "/at") || // VERY TEMPORARY
			strings.HasPrefix(c.Path(), "/static") ||
			strings.HasPrefix(c.Path(), "/favicon.ico") || //skip static assets
			strings.HasSuffix(c.Path(), "ws") //||
	}, jwtmiddleware.JWTBearerTokenExtractor, rolesRepository))

	e.Use(jwtmiddleware.JwtAuth(func(c echo.Context) bool { //websocket auth
		return !strings.HasSuffix(c.Path(), "ws")
	}, jwtmiddleware.JWTProtocolHeaderExtractor, rolesRepository))
}

func initDatabaseRepositories() (db.UserRepository, db.RolesRepository,
		db.RoomsRepository, db.MessagesRepository, db.VerificationCodeRepository){
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

	roomsRepository, err := rooms.New(database)
	if err != nil {
		panic(err)
	}

	messagesRepository, err := messages.New(database)
	if err != nil {
		panic(err)
	}

	verificationsRepo, err := verifications.New(database)
	if err != nil {
		panic(err)
	}

	return userRepository, rolesRepository, roomsRepository, messagesRepository, verificationsRepo
}

func main() {
	e := echo.New()
	e.Debug = config.Debug

	e.Static("/web", "frontend/build")
	e.Static("/static", "frontend/build/static")
	e.Static("/service-worker.js", "frontend/build/service-worker.js")

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

	userRepository, rolesRepository, roomsRepository,
		messagesRepository, verificationsRepository = initDatabaseRepositories()


	addMiddlewares(e, rolesRepository)
	createApiRoutes(v1ApiGroup, hubMap)

	t := &Template{
		templates: template.Must(template.ParseGlob("frontend/build/*.html")),
	}
	e.Renderer = t
	e.GET("/", controllers.Index)
	e.GET("/at/*", controllers.Index)
	e.GET("/wstest", controllers.WSTestView)

	//log.SetOutput(os.Stdout)
	e.Logger.Fatal(e.Start(":4000"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}