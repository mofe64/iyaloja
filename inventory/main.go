package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/config"
	"github.com/mofe64/iyaloja/inventory/middleware"
	"github.com/mofe64/iyaloja/inventory/route"
	"github.com/mofe64/iyaloja/inventory/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.New()
	router.Use(middleware.CustomLogger())
	router.Use(gin.Recovery())

	// Set up Database Connection
	config.ConnectDB()

	// Set up ping route
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	//Application Routes
	route.InventoryRoute(router)

	// Create Custom Server
	server := &http.Server{
		Addr:    ":" + config.EnvHTTPPort(),
		Handler: router,
	}

	go func() {
		util.ApplicationLog.Println("Starting server on port " + config.EnvHTTPPort())
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			util.ApplicationLog.Printf("error listen: %v\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	util.ApplicationLog.Println("Received signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		util.ApplicationLog.Fatal("Server forced to shutdown:", err)
	}

	util.ApplicationLog.Println("Server exiting....")

}
