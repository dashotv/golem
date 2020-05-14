// This file is autogenerated by Golem
// Do NOT make modifications, they will be lost
package torrents

import (
	"github.com/gin-gonic/gin"

	"github.com/dashotv/golem/web"
	"github.com/dashotv/test/config"
)

var cfg *config.Config

func Routes(c *config.Config, router *gin.Engine) {
	cfg = c
	r := router.Group("/torrents")
	r.GET("/add", addHandler)
	r.GET("/destroy", destroyHandler)
	r.GET("/", indexHandler)
	r.GET("/stop", labelHandler)
	r.GET("/pause", pauseHandler)
	r.GET("/remove", removeHandler)
	r.GET("/resume", resumeHandler)
	r.GET("/start", startHandler)
	r.GET("/stop", stopHandler)
	r.GET("/want", wantHandler)
	r.GET("/wanted", wantedHandler)

}

func addHandler(c *gin.Context) {
	url := web.QueryString(c, "url")

	Add(c, url)
}

func destroyHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Destroy(c, infohash)
}

func indexHandler(c *gin.Context) {

	Index(c)
}

func labelHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")
	label := web.QueryString(c, "label")

	Label(c, infohash, label)
}

func pauseHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Pause(c, infohash)
}

func removeHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Remove(c, infohash)
}

func resumeHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Resume(c, infohash)
}

func startHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Start(c, infohash)
}

func stopHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Stop(c, infohash)
}

func wantHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")
	files := web.QueryString(c, "files")

	Want(c, infohash, files)
}

func wantedHandler(c *gin.Context) {
	infohash := web.QueryString(c, "infohash")

	Wanted(c, infohash)
}
