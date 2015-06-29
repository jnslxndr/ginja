package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jens-a-e/ginja"
)

func ApiReporter(c *gin.Context) {
	log.Println("ginja middleware works!")
	c.Next()
}

func main() {
	server := gin.Default()

	api := ginja.New(server, ginja.Config{
		MountStats: true,
	}, ApiReporter)

	api.Register(Item{})

	server.Run(":8080")
}
