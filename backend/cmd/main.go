package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // change in PROD
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
	go func() {
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}
		log.Println("attempt to start the server")
		log.Println("Server bhiiii chaluuuuuuuuuuuuuuuuuu , READYY to serveeee ;)")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server not started %v", err)
			return
		}

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting DOWNWWW!!")
}
