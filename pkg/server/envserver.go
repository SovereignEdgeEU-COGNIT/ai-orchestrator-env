package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type EnvServer struct {
	ginHandler *gin.Engine
	port       int
	httpServer *http.Server
	db         database.Database
}

func CreateEnvServer(db database.Database, port int) *EnvServer {
	server := &EnvServer{}
	server.ginHandler = gin.Default()
	server.ginHandler.Use(cors.Default())

	server.db = db

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: server.ginHandler,
	}

	server.httpServer = httpServer
	server.port = port

	log.WithFields(log.Fields{"Port": port}).Info("Starting EnvServer")

	server.setupRoutes()

	return server
}

func (server *EnvServer) setupRoutes() {
	server.ginHandler.POST("/metrics", server.handleAddMetricRequest)
	server.ginHandler.GET("/metrics", server.handleGetMetricsRequest)
}

func (server *EnvServer) ServeForever() error {
	if err := server.httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (server *EnvServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{"Error": err}).Warning("EnvServer forced to shutdown")
	}
}
