package hostroute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HostConfig holds the configuration for each host.
type HostConfig struct {
	Host          string                 // The specific hostname (e.g., "host1.com").
	Prefix        string                 // Prefix for route paths, allowing access to specific routes on a generic host (e.g., "1" or "2").
	RouterFactory func(*gin.RouterGroup) // Function to define routes for the host.
	engine        *gin.Engine            // Internal engine configured for the host.
}

// createHostBasedRoutingMiddleware returns a middleware function to manage routes based on host configuration and secure against unknown hosts.
func createHostBasedRoutingMiddleware(hostConfigMap map[string]*HostConfig, genericHosts map[string]bool, secureAgainstUnknownHosts bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		host := c.Request.Host

		if config, exists := hostConfigMap[host]; exists {
			// Serve requests using the specific engine for the recognized host.
			config.engine.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		if _, exists := genericHosts[host]; exists {
			// Continue with the request processing for generic hosts.
			c.Next()
			return
		}

		if secureAgainstUnknownHosts {
			// If host is not recognized, return a 404 Not Found response.
			c.String(http.StatusNotFound, "Unknown host")
			c.Abort()
			return
		}

		// Proceed with the request processing if not secured against unknown hosts.
		c.Next()
	}
}

// SetupHostBasedRoutes configures routing based on hostnames using the given engine and configuration options.
func SetupHostBasedRoutes(r *gin.Engine, hostConfigs []HostConfig, genericHosts []string, secureAgainstUnknownHost bool, additionalHostConfig ...func(string, *gin.Engine) error) error {
	hostConfigMap := make(map[string]*HostConfig)
	genericHostsMap := stringSliceToMap(genericHosts)

	for i := range hostConfigs {
		engine := gin.New()        // Create a new Gin Engine for each host configuration.
		engine.Use(gin.Recovery()) // Add recovery middleware for handling panics.

		if len(additionalHostConfig) > 0 {
			// Apply any additional host configurations provided.
			for _, config := range additionalHostConfig {
				err := config(hostConfigs[i].Host, engine)
				if err != nil {
					return err // Return on error in configuration.
				}
			}
		}

		hostConfigs[i].engine = engine
		hostConfigs[i].RouterFactory(&engine.RouterGroup) // Set up routes using the provided factory function.

		if hostConfigs[i].Prefix != "" {
			// Create prefixed routes for generic hosts.
			group := r.Group(fmt.Sprintf("/%s", hostConfigs[i].Prefix))
			hostConfigs[i].RouterFactory(group) // Set up routes for each prefix.
		}

		hostConfigMap[hostConfigs[i].Host] = &hostConfigs[i] // Maintain a map of host configurations.
	}

	// Apply additional configurations to generic hosts if provided.
	if len(additionalHostConfig) > 0 {
		for _, genericHost := range genericHosts {
			for _, config := range additionalHostConfig {
				err := config(genericHost, r)
				if err != nil {
					return err // Return on setup error.
				}
			}
		}
	}

	// Apply middleware to manage host-based routing and security.
	r.Use(createHostBasedRoutingMiddleware(hostConfigMap, genericHostsMap, secureAgainstUnknownHost))

	return nil
}

// stringSliceToMap converts a slice of strings to a map with the string as the key and true as the value.
func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range slice {
		result[s] = true
	}
	return result
}
