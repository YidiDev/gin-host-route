package hostroute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func defineHost1Routes(g *gin.RouterGroup) {
	g.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host1")
	})
	g.GET("/hi", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi from host1")
	})
}

func defineHost2Routes(g *gin.RouterGroup) {
	g.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from host2")
	})
	g.GET("/hi", func(c *gin.Context) {
		c.String(http.StatusOK, "Hi from host2")
	})
}

func noRouteSpecifier(_ string, r *gin.Engine) error {
	r.NoRoute(noRouteHandler)
	return nil
}

func noRouteHandler(c *gin.Context) {
	c.String(http.StatusNotFound, "No known route")
}

func TestHostBasedRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.NoRoute(noRouteHandler)

	hostConfigs := []HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	err := SetupHostBasedRoutes(r, hostConfigs, genericHosts, true, noRouteSpecifier)
	require.NoError(t, err)

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		host, path string
		expected   string
		statusCode int
	}{
		{"host1.com", "/", "Hello from host1", http.StatusOK},
		{"host1.com", "/hi", "Hi from host1", http.StatusOK},
		{"host1.com", "/unknown", "No known route", http.StatusNotFound},

		{"host2.com", "/", "Hello from host2", http.StatusOK},
		{"host2.com", "/hi", "Hi from host2", http.StatusOK},
		{"host2.com", "/unknown", "No known route", http.StatusNotFound},

		{"host3.com", "/1", "Hello from host1", http.StatusOK},
		{"host3.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host3.com", "/2", "Hello from host2", http.StatusOK},
		{"host3.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host3.com", "/unknown", "No known route", http.StatusNotFound},

		{"host4.com", "/1", "Hello from host1", http.StatusOK},
		{"host4.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host4.com", "/2", "Hello from host2", http.StatusOK},
		{"host4.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host4.com", "/unknown", "No known route", http.StatusNotFound},

		{"unknown.com", "/", "Unknown host", http.StatusNotFound},
	}

	client := &http.Client{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Host: %s, Path: %s", tt.host, tt.path), func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL+tt.path, nil)
			req.Host = tt.host

			var resp *http.Response
			resp, err = client.Do(req)

			require.NoError(t, err)
			defer func() {
				err = resp.Body.Close()
				require.NoError(t, err)
			}()

			var body []byte
			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			assert.Equal(t, tt.expected, string(body))
		})
	}
}

func TestHostBasedRoutingWithoutSecure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.NoRoute(noRouteHandler)

	hostConfigs := []HostConfig{
		{Host: "host1.com", Prefix: "1", RouterFactory: defineHost1Routes},
		{Host: "host2.com", Prefix: "2", RouterFactory: defineHost2Routes},
	}

	genericHosts := []string{"host3.com", "host4.com"}

	err := SetupHostBasedRoutes(r, hostConfigs, genericHosts, false, noRouteSpecifier)
	require.NoError(t, err)

	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		host, path string
		expected   string
		statusCode int
	}{
		{"host1.com", "/", "Hello from host1", http.StatusOK},
		{"host1.com", "/hi", "Hi from host1", http.StatusOK},
		{"host1.com", "/unknown", "No known route", http.StatusNotFound},

		{"host2.com", "/", "Hello from host2", http.StatusOK},
		{"host2.com", "/hi", "Hi from host2", http.StatusOK},
		{"host2.com", "/unknown", "No known route", http.StatusNotFound},

		{"host3.com", "/1", "Hello from host1", http.StatusOK},
		{"host3.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host3.com", "/2", "Hello from host2", http.StatusOK},
		{"host3.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host3.com", "/unknown", "No known route", http.StatusNotFound},

		{"host4.com", "/1", "Hello from host1", http.StatusOK},
		{"host4.com", "/1/hi", "Hi from host1", http.StatusOK},
		{"host4.com", "/2", "Hello from host2", http.StatusOK},
		{"host4.com", "/2/hi", "Hi from host2", http.StatusOK},
		{"host4.com", "/unknown", "No known route", http.StatusNotFound},

		{"unknown.com", "/", "No known route", http.StatusNotFound},
	}

	client := &http.Client{}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Host: %s, Path: %s", tt.host, tt.path), func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL+tt.path, nil)
			req.Host = tt.host

			var resp *http.Response
			resp, err = client.Do(req)

			require.NoError(t, err)
			defer func() {
				err = resp.Body.Close()
				require.NoError(t, err)
			}()

			var body []byte
			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, resp.StatusCode)
			assert.Equal(t, tt.expected, string(body))
		})
	}
}
