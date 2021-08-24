package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	httpserver "github.com/daqnext/meson.network-lts-http-server"
)

func main() {
	hs := httpserver.New()
	hs.SetLogLevel_DEBUG()

	//////////start transmiting after 10 seconds////////
	hs.StaticWithPause("/", "assets")
	hs.SetPauseSeconds(10)
	///////////log related setting/////////////////////////
	fmt.Println("log level is :", hs.Logger.Level())
	hs.SetLogFile("./output2.txt")

	///////////////////  JSONP //////////////////////
	hs.GET("/jsonp", func(c httpserver.Context) error {
		callback := c.QueryParam("callback")
		var content struct {
			Response  string    `json:"response"`
			Timestamp time.Time `json:"timestamp"`
			Random    int       `json:"random"`
		}
		content.Response = "Sent via JSONP"
		content.Timestamp = time.Now().UTC()
		content.Random = rand.Intn(1000)
		return c.JSONP(http.StatusOK, callback, &content)
	})

	hs.GET("/sendfiletest/:filename", func(c httpserver.Context) error {
		name := c.Param("filename")
		needSavedHeader := true
		return httpserver.FileWithPause(hs, c, "assets/"+name, needSavedHeader)
	})

	///////////////////start//////////////////////////////
	hs.Logger.Fatal(hs.Start(":8080")) //stuck here

	//////////////////start using https server////////////
	//hs.Logger.Fatal(hs.StartTLS(......))

	//////////////////start using https server////////////
	/////hs.Logger.Fatal(hs.StartAutoTLS(.....))

	////////don't forget to realse finally//////////////////
	//hs.CloseServer() //don't forget to closeserver

}
