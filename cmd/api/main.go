package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-crud-employee/config"
	"go-crud-employee/database"
	"go-crud-employee/handlers"
	"go-crud-employee/middleware"
	"go-crud-employee/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create database tables
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	// Create default admin user
	if err := db.CreateDefaultUser(); err != nil {
		log.Printf("Warning: Failed to create default user: %v", err)
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, jwtManager)
	employeeHandler := handlers.NewEmployeeHandler(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Setup router
	router := setupRouter(authHandler, employeeHandler, authMiddleware)

	// Start server
	log.Printf("Starting server on %s", cfg.GetServerAddress())
	log.Printf("Environment: %s", cfg.Env)
	log.Printf("Default admin credentials - Username: admin, Password: admin123")

	// Graceful shutdown
	go func() {
		if err := router.Run(cfg.GetServerAddress()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

func setupRouter(authHandler *handlers.AuthHandler, employeeHandler *handlers.EmployeeHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Employee Management API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)
		}

		// Employee routes (protected)
		employees := v1.Group("/employees")
		employees.Use(authMiddleware.RequireAuth())
		{
			employees.POST("/", employeeHandler.CreateEmployee)
			employees.GET("/", employeeHandler.GetEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}
	}

	return router
}