package main

import httpserver "github.com/daqnext/meson.network-lts-http-server"

func main() {
	hs := httpserver.New()
	hs.StaticWithPause(hs, "/", "assets")
	hs.Logger.Fatal(hs.Start(":80"))

}
