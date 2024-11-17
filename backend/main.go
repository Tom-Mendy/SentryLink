package main

import (
	"net/http"
	"os"

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
	"github.com/Tom-Mendy/SentryLink/service"
)

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
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	var (
		// Database connection
		databaseConnection *gorm.DB = database.Connection()

		// Repositories
		linkRepository        repository.LinkRepository        = repository.NewLinkRepository(databaseConnection)
		githubTokenRepository repository.GithubTokenRepository = repository.NewGithubTokenRepository(databaseConnection)
		userRepository        repository.UserRepository        = repository.NewUserRepository(databaseConnection)

		// Services
		linkService        service.LinkService        = service.NewLinkService(linkRepository)
		githubTokenService service.GithubTokenService = service.NewGithubTokenService(githubTokenRepository)
		userService        service.UserService        = service.NewUserService(userRepository)
		jwtService         service.JWTService         = service.NewJWTService()

		// Controllers
		linkController        controller.LinkController        = controller.NewLinkController(linkService)
		githubTokenController controller.GithubTokenController = controller.NewGithubTokenController(githubTokenService)
		userController        controller.UserController        = controller.NewUserController(userService, jwtService)
	)

	linkApi := api.NewLinkAPI(userController, linkController, githubTokenController)

	apiRoutes := router.Group(docs.SwaggerInfo.BasePath)
	{
		// User Auth
		auth := apiRoutes.Group("/auth")
		{
			auth.POST("/login", linkApi.Login)
			auth.POST("/register", linkApi.Register)
		}

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
				linkApi.RedirectToGithub(c, github.BasePath()+"/auth/callback")
			})

			github.GET("/auth/callback", func(c *gin.Context) {
				linkApi.HandleGithubTokenCallback(c, github.BasePath()+"/auth/callback")
			})
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// view request received but not found
	router.NoRoute(func(c *gin.Context) {
		// get the path
		path := c.Request.URL.Path
		// get the method
		method := c.Request.Method
		print("\n\n" + method + " " + path + "\n\n\n")
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
	router := setupRouter()

	// Listen and Server in 0.0.0.0:8000
	err := router.Run(":8000")
	if err != nil {
		panic("Error when running the server")
	}
}
