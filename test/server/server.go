package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dashotv/blarg/config"
)

type Server struct {
	Config *config.Config
	Router *gin.Engine
	Log    *logrus.Entry
	// additional server needs, e.g. cache client, database client, etc
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{Config: cfg}

	if cfg.Mode == "release" {
		gin.SetMode(cfg.Mode)
	}

	host, _ := os.Hostname()
	s.Log = logrus.WithField("prefix", host)
	s.Log.Level = logrus.DebugLevel

	return s, nil
}

func (s *Server) Start() error {
	s.Log.Info("starting Blarg...")

	s.Router = gin.Default()
	s.Routes()

	//s.Jobs configuration

	s.Log.Info("starting web...")
	if err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port)); err != nil {
		return errors.Wrap(err, "starting router")
	}

	return nil
}

func (s *Server) Routes() {
	s.Router.GET("/", homeIndex)

}

func homeIndex(c *gin.Context) {
	c.String(http.StatusOK, "home")
}
