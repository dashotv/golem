package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/test/server/nzbs"
	"github.com/dashotv/test/server/torrents"
)

func (s *Server) Routes() {
	s.Router.GET("/", homeHandler)

	nzbs.Routes(s.Config, s.Router)
	torrents.Routes(s.Config, s.Router)

}

func homeHandler(c *gin.Context) {
	Home(c)
}

func Home(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
