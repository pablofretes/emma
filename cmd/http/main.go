package main

import (
	"emma/cmd/http/docs"
	"fmt"
	"log"
	"net/http"

	// docs "github.com/emma/cmd/http/docs"
	configs "emma/cmd/configs"
	routers "emma/cmd/http/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/docs/index.html")
	})

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