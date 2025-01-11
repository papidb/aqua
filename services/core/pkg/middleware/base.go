package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/papidb/aqua/pkg/config"
	"github.com/rs/zerolog"
)

// PrepareRequest hooks fundamental utilities (cors, compression, source IP, logging, recoveries and timeouts)
// into the request context. Ensure the to Use this function before the request hits any handler
func PrepareRequest(app *config.App, r *gin.Engine, log zerolog.Logger) *gin.Engine {
	r.Use(enableCors(app))
	r.Use(gin.Recovery())
	return r
}
