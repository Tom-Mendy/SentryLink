package main

import (
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Tom-Mendy/SentryLink/api"
	"github.com/Tom-Mendy/SentryLink/controller"
	"github.com/Tom-Mendy/SentryLink/database"
	"github.com/Tom-Mendy/SentryLink/docs"
	"github.com/Tom-Mendy/SentryLink/middlewares"
	"github.com/Tom-Mendy/SentryLink/repository"
	"github.com/Tom-Mendy/SentryLink/schemas"
	"github.com/Tom-Mendy/SentryLink/service"
)

type ActionService struct {
	Service string
	Action  string
}

func hello(c chan ActionService) {
	var dt time.Time
	for {
		dt = time.Now().Local()
		if dt.Hour() == 12 && dt.Minute() == 56 {
			println("current time is ", dt.String())
			c <- ActionService{
				Service: "Timer",
				Action:  "say Hello",
			} // send sum to c
		}
		time.Sleep(30 * time.Second)
	}
}

func world(c chan ActionService) {
	var dt time.Time
	for {
		dt = time.Now().Local()
		if dt.Hour() == 12 && dt.Minute() == 56 {
			println("current time is ", dt.String())
			c <- ActionService{
				Service: "Timer",
				Action:  "say World",
			} // send sum to c
		}
		time.Sleep(20 * time.Second)
	}
}

func setupRouter() *gin.Engine {

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		panic("APP_PORT is not set")
	}

	docs.SwaggerInfo.Title = "SentryLink API"
	docs.SwaggerInfo.Description = "SentryLink - Crawler API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + appPort
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router := gin.Default()

	// Ping test
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, &schemas.Response{
			Message: "pong",
		})
	})

	var (
		// Database connection
		databaseConnection *gorm.DB = database.Connection()

		// Repositories
		linkRepository        repository.LinkRepository        = repository.NewLinkRepository(databaseConnection)
		githubTokenRepository repository.GithubTokenRepository = repository.NewGithubTokenRepository(databaseConnection)
		userRepository        repository.UserRepository        = repository.NewUserRepository(databaseConnection)

		// Services
		jwtService         service.JWTService         = service.NewJWTService()
		linkService        service.LinkService        = service.NewLinkService(linkRepository)
		githubTokenService service.GithubTokenService = service.NewGithubTokenService(githubTokenRepository)
		userService        service.UserService        = service.NewUserService(userRepository, jwtService)

		// Controllers
		linkController        controller.LinkController        = controller.NewLinkController(linkService)
		githubTokenController controller.GithubTokenController = controller.NewGithubTokenController(githubTokenService, userService)
		userController        controller.UserController        = controller.NewUserController(userService, jwtService)
	)

	linkApi := api.NewLinkAPI(linkController)

	userApi := api.NewUserAPI(userController)

	githubApi := api.NewGithubAPI(githubTokenController)

	apiRoutes := router.Group(docs.SwaggerInfo.BasePath)
	{
		// User Auth
		auth := apiRoutes.Group("/auth")
		{
			auth.POST("/login", userApi.Login)
			auth.POST("/register", userApi.Register)
		}

		// Links
		links := apiRoutes.Group("/links", middlewares.AuthorizeJWT())
		{
			links.GET("", linkApi.GetLink)
			links.POST("", linkApi.CreateLink)
			links.PUT(":id", linkApi.UpdateLink)
			links.DELETE(":id", linkApi.DeleteLink)
		}

		// Github
		github := apiRoutes.Group("/github")
		{
			github.GET("/auth", func(c *gin.Context) {
				githubApi.RedirectToGithub(c, github.BasePath()+"/auth/callback")
			})

			github.GET("/auth/callback", func(c *gin.Context) {
				githubApi.HandleGithubTokenCallback(c, github.BasePath()+"/auth/callback")
			})

			githubInfo := github.Group("/info", middlewares.AuthorizeJWT())
			{
				githubInfo.GET("/user", githubApi.GetUserInfo)
			}

		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// view request received but not found
	router.NoRoute(func(c *gin.Context) {
		// get the path
		path := c.Request.URL.Path
		// get the method
		method := c.Request.Method
		c.JSON(http.StatusNotFound, gin.H{"error": "not found", "path": path, "method": method})
	})

	return router
}

func init() {
	// err := .Load()
	// if err != nil {
	// 	panic("Error loading .env file")
	// }
}

// @securityDefinitions.apiKey bearerAuth
// @in header
// @name Authorization
func main() {

	// Create a channel list
	var mychannel1 = make([]chan ActionService, 2)
	mychannel1[0] = make(chan ActionService)
	mychannel1[1] = make(chan ActionService)
	mychannel2 := make(chan ActionService)
	mychannel1 = append(mychannel1, mychannel2)

	go hello(mychannel1[0])
	go world(mychannel1[0])
	go world(mychannel1[2])

	go func(mychannel1 []chan ActionService) {
		for {
			x := <-mychannel1[0]
			if x.Service == "Timer" {
				println(x.Action)
			} else {
				println("Unknown service")
			}
		}
	}(mychannel1)

	router := setupRouter()

	// Listen and Server in 0.0.0.0:8000
	err := router.Run(":8080")
	if err != nil {
		panic("Error when running the server")
	}
}
