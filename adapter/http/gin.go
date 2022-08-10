package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type GinHandler struct {
	engine *gin.Engine
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func NewGinHandler() *GinHandler {
	gin.SetMode(gin.ReleaseMode)
	return &GinHandler{
		engine: gin.New(),
	}
}

func (f *GinHandler) Run(address string) {
	log.Infof("Listening on %s", address)
	f.engine.Run(address)
}

func (f *GinHandler) SetupRoutes() {
	f.engine.GET("/metrics", prometheusHandler())
}
