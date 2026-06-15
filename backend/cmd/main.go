package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	//"github.com/PrasadNaik1310/LMSVR_SM/handlers"
	"github.com/PrasadNaik1310/LMSVR_SM/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting new server")

	if err := godotenv.Load(); err != nil {
		log.Println("Env Env file not found , main file error")

	}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		if os.Getenv("APP_ENV") == "dev" {
			log.Printf("WARNING: Dev environment configured , ALLOWING ALL ORIGINS")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("allowed_origin"))
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "PUT,POST,OPTIONS,GET,DELETE,PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			log.Println("Recived preflight request !!")
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	err := db.InitDB()
	if err != nil {
		log.Println(err)
		log.Fatalf("Error init. Db , from main")
		return
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Service": "LMSVR_SM",
			"Status":  "Healthy",
		})

	})
	// all routes are registered and managed in routes folder. No route is being called here in the main file.
	routes.RegisterRoutes(r)
	port := os.Getenv("port")
	if port == "" {
		log.Println("Port not found ")
		port = "8080"
	}
	/*	if os.Getenv("APP_ENV") != "production" {
		if err := seed.MigrateAndSeed(db); err != nil {
			log.Fatalf(err)
			return
		}
	}*/
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {

		log.Println("Starting up server in gorountine: Complete!!")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server not started %v", err)
			return
		}

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("signal for  gracefull shutdown: attempt started")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("CoudlNot perform graceful shutdown")
		log.Fatalf("Server FORCED to shutdown : timelimit excedded for graceful shutdown")
	}
	fmt.Println("Signal for graceful shutdown : Completed!!")
}
