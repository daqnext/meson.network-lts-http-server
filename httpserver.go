package MSHttpServer

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Context = echo.Context
type HttpServer struct {
	echo.Echo
}

func (e *HttpServer) StaticWithPause(prefix, root string) *echo.Route {
	if root == "" {
		root = "." // For security we want to restrict to CWD.
	}
	return e.static_with_pause(prefix, root, e.GET)
}

func (e *HttpServer) static_with_pause(prefix, root string, get func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route) *echo.Route {
	h := func(c Context) error {
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}

		name := filepath.Join(root, filepath.Clean("/"+p)) // "/"+ for security
		fi, err := os.Stat(name)
		if err != nil {
			// The access path does not exist
			return echo.NotFoundHandler(c)
		}

		// If the request is for a directory and does not end with "/"
		p = c.Request().URL.Path // path must not be empty.
		if fi.IsDir() && p[len(p)-1] != '/' {
			// Redirect to ends with "/"
			return c.Redirect(http.StatusMovedPermanently, p+"/")
		}
		return c.File(name)
	}
	// Handle added routes based on trailing slash:
	// 	/prefix  => exact route "/prefix" + any route "/prefix/*"
	// 	/prefix/ => only any route "/prefix/*"
	if prefix != "" {
		if prefix[len(prefix)-1] == '/' {
			// Only add any route for intentional trailing slash
			return get(prefix+"*", h)
		}
		get(prefix, h)
	}
	return get(prefix+"/*", h)
}

func New() (hs *HttpServer) {
	hs = &HttpServer{*echo.New()}
	return hs
}
