# Host Route Library

A high-performance Gin middleware library for routing based on the host. This library facilitates the configuration of different routes and behaviors for distinct hostnames, enhancing the ability to host multi-tenant applications on a single server.

[![go report card](https://goreportcard.com/badge/github.com/YidiDev/gin-host-route "go report card")](https://goreportcard.com/report/github.com/YidiDev/gin-host-route)
[![test status](https://github.com/YidiDev/gin-host-route/workflows/tests/badge.svg?branch=main "test status")](https://github.com/YidiDev/gin-host-route/actions)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/YidiDev/gin-host-route?tab=doc)

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

### Additional Host Config
This param is optional and allows for unlimited inputs. Each input should be a `func(*echo.Group) error`. This is meant for specifying functions that `SetupHostBasedRoutes` should run on every host group after creating it. Common use cases of this are:
- Configuring a `NoRoute` Handler.
- Configuring Host Specific Middleware. This can be done in the `HostConfig` in the `RouterFactory`. Alternatively, it could be done here. This may be useful if you want to centralize a lot of the host-specific middleware.

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

## Sister Project
This project has a sister project for Echo framework users. If you are using Echo, check out the [Echo Host Route Library](https://github.com/YidiDev/echo-host-route) for similar functionality.

## Contributing
Contributions are always welcome! If you're interested in contributing to the project, please take a look at our [Contributing Guidelines](CONTRIBUTING.md) file for guidelines on how to get started. We appreciate your help in improving the library!
