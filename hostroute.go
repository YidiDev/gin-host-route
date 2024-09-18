package hostroute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HostConfig struct {
	Host          string
	Prefix        string
	RouterFactory func(*gin.RouterGroup)
	engine        *gin.Engine
}

func createHostBasedRoutingMiddleware(hostConfigMap map[string]*HostConfig, genericHosts map[string]bool, secureAgainstUnknownHosts bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		host := c.Request.Host

		if config, exists := hostConfigMap[host]; exists {
			config.engine.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		if _, exists := genericHosts[host]; exists {
			c.Next()
			return
		}

		if secureAgainstUnknownHosts {
			c.String(http.StatusNotFound, "Unknown host")
			c.Abort()
			return
		}

		c.Next()
	}
}

func SetupHostBasedRoutes(r *gin.Engine, hostConfigs []HostConfig, genericHosts []string, noRouteFactory func(*gin.Engine), secureAgainstUnknownHost bool) {
	hostConfigMap := make(map[string]*HostConfig)
	genericHostsMap := stringSliceToMap(genericHosts)

	for i := range hostConfigs {
		engine := gin.New()
		engine.Use(gin.Recovery())
		noRouteFactory(engine)
		hostConfigs[i].engine = engine
		hostConfigs[i].RouterFactory(&engine.RouterGroup)

		if hostConfigs[i].Prefix != "" {
			group := r.Group(fmt.Sprintf("/%s", hostConfigs[i].Prefix))
			hostConfigs[i].RouterFactory(group)
		}

		hostConfigMap[hostConfigs[i].Host] = &hostConfigs[i]
	}

	r.Use(createHostBasedRoutingMiddleware(hostConfigMap, genericHostsMap, secureAgainstUnknownHost))
}

func stringSliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, s := range slice {
		result[s] = true
	}
	return result
}
