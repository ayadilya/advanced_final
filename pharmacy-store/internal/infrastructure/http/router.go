package http

import (
	"database/sql"
	"net/http"
	"pharmacy-store/internal/handlers"
	natsClient "pharmacy-store/internal/infrastructure/messaging/nats"
	"pharmacy-store/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, natsClient *natsClient.NatsClient) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger())

	productHandler := handlers.NewProductHandler(db, natsClient)
	userHandler := handlers.NewUserHandler(db, natsClient)
	categoryHandler := handlers.NewCategoryHandler(db)

	api := router.Group("/api")
	{
		api.POST("/users/login", userHandler.Login)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/user/info", userHandler.GetUserInfo) // Add this line for user info

		secured := api.Group("/")
		secured.Use(middleware.JWTAuthMiddleware())
		{
			secured.GET("/products", productHandler.GetProducts)
			secured.GET("/products/:id", productHandler.GetProduct)
			secured.POST("/products", productHandler.CreateProduct)
			secured.PUT("/products/:id", productHandler.UpdateProduct)
			secured.DELETE("/products/:id", productHandler.DeleteProduct)

			secured.GET("/users", userHandler.GetUsers)
			secured.GET("/users/:id", userHandler.GetUser)
			secured.PUT("/users/:id", userHandler.UpdateUser)
			secured.DELETE("/users/:id", userHandler.DeleteUser)

			// Category routes
			secured.GET("/categories", categoryHandler.GetCategories)
			secured.GET("/categories/:id", categoryHandler.GetCategory)
			secured.POST("/categories", categoryHandler.CreateCategory)
			secured.PUT("/categories/:id", categoryHandler.UpdateCategory)
			secured.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			// Test NATS Publish Route
			secured.POST("/test-nats", func(c *gin.Context) {
				err := natsClient.Conn.Publish("test.subject", []byte("Hello from /test-nats"))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message", "details": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Message published"})
			})
		}
	}

	return router
}
