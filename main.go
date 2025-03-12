package main

import (
	"fmt"
	"os"

	"github.com/EputraP/kfc_be/internal/handler"
	"github.com/EputraP/kfc_be/internal/middleware"
	"github.com/EputraP/kfc_be/internal/repository"
	"github.com/EputraP/kfc_be/internal/routes"
	"github.com/EputraP/kfc_be/internal/service"
	dbstore "github.com/EputraP/kfc_be/internal/store"
	"github.com/EputraP/kfc_be/internal/util/hasher"
	"github.com/EputraP/kfc_be/internal/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func main() {

	if err := logger.Init("app.log"); err != nil {
		logger.Error("main", "Failed to initialize logger:", map[string]string{
			"error": err.Error(),
		})
		return
	}

	logger.Info("main", "Starting application...", nil)

	err := godotenv.Load()
	if err != nil {
		logger.Error("main", "Error loading .env file", map[string]string{
			"error": err.Error(),
		})
		return
	}

	handlers, middlewares := prepare()

	srv := gin.Default()

	srv.Use(middleware.CORS())

	routes.Build(srv, handlers, middlewares)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	logger.Info("main",
		"Server is starting...", map[string]string{
			"port": port,
		})

	if err := srv.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error("main", "Error running gin server:", map[string]string{
			"error": err.Error(),
		})
	}

}

func prepare() (handlers *routes.Handlers, middleware *routes.Middlewares) {
	logger.Info("main", "Initializing dependencies...", nil)
	hasher := hasher.NewBcrypt(10)
	db := dbstore.Get()

	logger.Info("main", "Initializing repositories...", nil)
	authRepo := repository.NewAuthRepository(db)

	logger.Info("main", "Initializing services...", nil)
	authService := service.NewAuthService(service.AuthServiceConfig{AuthRepo: authRepo, Hasher: hasher})

	logger.Info("main", "Initializing handlers...", nil)
	authHandler := handler.NewAuthHandler(handler.AuthHandlerConfig{AuthService: authService})

	handlers = &routes.Handlers{
		Auth: authHandler,
	}

	logger.Info("main", "Application initialized successfully.", nil)
	return
}
