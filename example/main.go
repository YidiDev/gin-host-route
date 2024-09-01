package main

import (
	"github.com/YidiDev/gin-host-route"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func defineHost1Routes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host1")
	})
	rg.GET("/hi", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi from host1")
	})
}

func defineHost2Routes(rg *gin.RouterGroup) {
	rg.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host2")
	})
	rg.GET("/hi", func(c *gin.Context) {
		log.Println("Important stuff")
		c.String(http.StatusOK, "Hi from host2")
	})
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	r := gin.Default()

	hostConfigs := []hostroute.HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	hostroute.SetupHostBasedRoutes(r, hostConfigs, genericHosts, true)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "No known route")
	})

	r.Run(":8080")
}
