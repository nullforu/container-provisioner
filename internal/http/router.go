package http

import (
	"io"
	nethttp "net/http"
	"os"

	"smctf/internal/config"
	"smctf/internal/http/handlers"
	"smctf/internal/http/middleware"
	"smctf/internal/logging"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, logger *logging.Logger) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if logger != nil {
		gin.DefaultWriter = io.MultiWriter(os.Stdout, logger)
		gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, logger)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(cfg.Logging, logger))

	_ = handlers.New(cfg)

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(nethttp.StatusOK, gin.H{"status": "ok"})
	})

	// api := r.Group("/api")

	return r
}
