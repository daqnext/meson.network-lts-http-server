package httpserver

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Context = echo.Context
type HttpServer struct {
	echo.Echo
	PauseMoment int64
	filehandler *os.File
}

func (hs *HttpServer) SetPauseSeconds(secs int64) {
	hs.PauseMoment = time.Now().Unix() + secs
}

func (hs *HttpServer) GetPauseMoment() int64 {
	return hs.PauseMoment
}

func FileWithPause(hs *HttpServer, c Context, file string, needSavedHeader bool) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, "index.html")
		f, err = os.Open(file)
		if err != nil {
			return echo.NotFoundHandler(c)
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			return
		}
	}
	ServeContent(hs, c.Response(), c.Request(), fi.Name(), fi.ModTime(), f, needSavedHeader)
	return
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

		return FileWithPause(e, c, name, true)
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
	hs = &HttpServer{*echo.New(), 0, nil}
	hs.Use(middleware.Logger())
	return hs
}

func (hs *HttpServer) SetLogOutput(writer io.Writer) error {
	hs.Logger.SetOutput(writer)
	return nil
}

func (hs *HttpServer) SetLogFile(path string) error {
	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	hs.filehandler = fh
	hs.Logger.SetOutput(hs.filehandler)
	return err
}

func (hs *HttpServer) CloseServer() {
	hs.Close()
	if hs.filehandler != nil {
		hs.filehandler.Close()
	}
}

func (hs *HttpServer) SetLogLevel_DEBUG() {
	hs.Logger.SetLevel(log.DEBUG)
}

func (hs *HttpServer) SetLogLevel_INFO() {
	hs.Logger.SetLevel(log.INFO)
}

func (hs *HttpServer) SetLogLevel_WARN() {
	hs.Logger.SetLevel(log.WARN)
}

func (hs *HttpServer) SetLogLevel_ERROR() {
	hs.Logger.SetLevel(log.ERROR)
}

func (hs *HttpServer) SetLogLevel_OFF() {
	hs.Logger.SetLevel(log.OFF)
}
