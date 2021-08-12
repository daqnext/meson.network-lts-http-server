package test

import (
	"testing"

	HttpServer "github.com/daqnext/meson.network-lts-http-server"
)

func TestMain(t *testing.T) {

	hs := HttpServer.New()
	hs.Static("/static", "/assets")
	hs.Start("80")
}
