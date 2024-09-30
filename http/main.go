package main

import (
	"fmt"
	"log"
	"net/http"

	configs "emma/configs"
	routers "emma/http/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Recovery())

	routers.InitRouters(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})

	fmt.Println("🌥️  Listen on port ", configs.GetConfig().PORT)
	fmt.Println("📘  Docs on port ", configs.GetConfig().PORT+"/docs")

	err := r.Run(configs.GetConfig().PORT)
	if err != nil {
		log.Println("Error running server", err)
		log.Fatal(err)
	}
}