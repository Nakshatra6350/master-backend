package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(os.Stdout, "Master :: ", log.LstdFlags)

	l.SetOutput(file)

	// Initialize Gin router
	router := setupRouter(l)
	server := &http.Server{
		Addr:         ":7005",
		Handler:      router,
		ReadTimeout:  60 * 6 * time.Second, // Maximum duration for reading the entire request
		WriteTimeout: 60 * 6 * time.Second, // Maximum duration before timing out writes of the response
		IdleTimeout:  120 * time.Second,    // Maximum duration an idle (keep-alive) connection will be kept open
		ErrorLog:     l,
	}

	// Start HTTP server in a separate goroutine
	func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(server)
}

func setupRouter(l *log.Logger) *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "Token", "client_id", "client_secret", "admin_secret"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	utility := router.Group("/utility")
	{
		hh := handlers.NewUtility(l)
		utility.GET("/health", hh.HealthCheckUp)
		utility.POST("/sendOtp", hh.SendOTP)
		utility.POST("/verifyOtp", hh.VerifyOTP)
	}
	users := router.Group("/user")
	{
		auth := handlers.NewUser(l)
		users.POST("/signup", auth.SignupHandler)
		users.POST("/login", auth.LoginHandler)
		users.POST("/logout", auth.LogoutHandler)
	}

	return router
}
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	// Set a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %s", err)
	}

	log.Println("Server gracefully stopped")
}
