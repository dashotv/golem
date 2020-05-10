package nzbs

import (
	"github.com/gin-gonic/gin"

	"github.com/dashotv/golem/web"
	"github.com/dashotv/test/server"
)

var serv *server.Server

func Routes(s *server.Server) {
	serv = s
	r := s.Router.Group("/nzbs")
	r.Get("/", indexHandler)
	r.Get("/add", addHandler)
	r.Get("/remove", removeHandler)
	r.Get("/destroy", destroyHandler)
	r.Get("/pause", pauseHandler)
	r.Get("/resume", resumeHandler)
	r.Get("/history", historyHandler)

}

func indexHandler(c *gin.Context) {

	index()
}

func addHandler(c *gin.Context) {
	url := web.QueryString(c, "url")
	category := web.QueryString(c, "category")
	name := web.QueryString(c, "name")

	add(url, category, name)
}

func removeHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	remove(id)
}

func destroyHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	destroy(id)
}

func pauseHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	pause(id)
}

func resumeHandler(c *gin.Context) {
	id := web.QueryInt(c, "id")

	resume(id)
}

func historyHandler(c *gin.Context) {
	hidden := web.QueryBool(c, "hidden")

	history(hidden)
}
