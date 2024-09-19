package main

import (
	"github.com/YidiDev/gin-host-route"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// defineHost1Routes sets up the routes specific to host1.com.
func defineHost1Routes(rg *gin.RouterGroup) {
	// Route to handle GET request to the root URL of host1, returns a greeting message.
	rg.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host1")
	})
	// Route to handle GET request to /hi URL of host1, returns a different greeting message.
	rg.GET("/hi", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi from host1")
	})
}

// defineHost2Routes sets up the routes specific to host2.com.
func defineHost2Routes(rg *gin.RouterGroup) {
	// Route to handle GET request to the root URL of host2, returns a greeting message.
	rg.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host2")
	})
	// Route to handle GET request to /hi URL of host2, includes log statement and returns a greeting message.
	rg.GET("/hi", func(c *gin.Context) {
		log.Println("Important stuff")
		c.String(http.StatusOK, "Hi from host2")
	})
}

// init function sets up logging output to standard output.
func init() {
	log.SetOutput(os.Stdout)
}

// noRouteHandler defines behavior for unspecified routes, returning a not-found message.
func noRouteHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "No known route")
}

// noRouteSpecifier applies a no-route handler for the Gin Engine.
func noRouteSpecifier(_ string, r *gin.Engine) error {
	r.NoRoute(noRouteHandler)
	return nil
}

// main function initializes the Gin engine and sets up host-based routing.
func main() {
	r := gin.Default() // Create a new Gin Engine instance with default middleware.

	// Define host-specific configurations using HostConfig structs.
	hostConfigs := []hostroute.HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	// Define generic hosts to use the primary router without specialized sub-routes.
	genericHosts := []string{"host3.com", "host4.com"}

	// Set up host-based routes and handle any setup errors.
	err := hostroute.SetupHostBasedRoutes(r, hostConfigs, genericHosts, true, noRouteSpecifier)
	if err != nil {
		log.Fatal(err) // Log fatal error if setup fails.
	}

	// Start the server on port 8080.
	r.Run(":8080")
}
