# Host Route Library

A high-performance Gin middleware library for routing based on the host.

## Installation

Add the module to your project by running:

```sh
go get github.com/YidiDev/gin-host-route
```

## Usage

Below is an example of how to utilize the library to define different routes based on the host.

### Example

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/YidiDev/gin-host-route"
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
		c.String(http.StatusOK, "Hi from host2")
	})
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	r := gin.Default()

	// Define host-specific configurations
	hostConfigs := []hostroute.HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	// Generic hosts are hosts that will use the primary router without special sub-routes
	genericHosts := []string{"host3.com", "host4.com"}

	// Setup host-based routes
	hostroute.SetupHostBasedRoutes(r, hostConfigs, genericHosts, true)

	// Define handler for unmatched routes
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "No known route")
	})

	// Start the server
	r.Run(":8080")
}
```

## Configuration Options

### `HostConfig`
The `HostConfig` struct is used to define the configuration for a specific host:
- `Host`: The hostname for which the configuration is defined.
- `Prefix`: A prefix to use for routes specific to this host when accessed on a generic host.
- `RouterFactory` A function that defined the routes for this host.

### Generic Hosts
Generic hosts are hosts that will have access to all routes defined in all the host configs and any others defined on the default router. This is useful for:
- **Local Testing**: to be able to access all routes without changing the host. 
- **Consolidated Access**: Handle routes from multiple applications on a single host. For example:
  - You have two applications hosted on one Go server: one at `application1.example.com` and the other at `application2.example.com`. However, you also want people to be able to access both applications by going to `example.com/application1` or `example.com/application2`.

### Secure Against Unknown Hosts
The `secureAgainstUnknownHosts` boolean flag controls how the middleware handles requests from unknown hosts:
- `true`: Requests from unknown hosts will receive a 404 Not Found Response. This is useful for securing your application against unexpected or unauthorized hosts.
- `false`: Requests from unknown hosts will be passed through the primary router. This is useful if you want to catch and handle such requests manually.

### Route Configuration Example

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/YidiDev/gin-host-route"
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
```

### Handling Different Hosts

1. **Host-specific Routes**:
   Routes are defined uniquely for each host using a specific `RouterFactory`. The `HostConfig` struct includes the hostname, path prefix, and a function to define routes for that host.

    ```go
    hostConfigs := []hostroute.HostConfig{
        {Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
        {Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
    }
    ```

2. **Generic Hosts**:
   Generic hosts allow for fallback to common routes defined in the primary router.

    ```go
    genericHosts := []string{"host3.com", "host4.com"}
    ```

3. **Secure Against Unknown Hosts**:
   Secure your application by handling unknown hosts, preventing them from accessing unintended routes.

    ```go
    hostroute.SetupHostBasedRoutes(r, hostConfigs, genericHosts, true)
    ```
