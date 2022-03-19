package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mofe64/iyaloja/inventory/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	applicationLog := log.New(os.Stdout, "inventory-service", log.LstdFlags)
	router := gin.New()
	router.Use(middleware.CustomLogger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		applicationLog.Println("Starting server on ...")
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			applicationLog.Printf("error listen: %v\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	applicationLog.Println("Received signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		applicationLog.Fatal("Server forced to shutdown:", err)
	}

	applicationLog.Println("Server exiting....")

}
