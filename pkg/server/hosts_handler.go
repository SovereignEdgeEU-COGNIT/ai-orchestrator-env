package server

import (
	"io"
	"net/http"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/gin-gonic/gin"
)

func (server *EnvServer) handleAddHostRequest(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	host, err := core.ConvertJSONToHost(string(jsonData))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = server.db.AddHost(host)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "")
}

func (server *EnvServer) handleGetHostRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "Paramater id must be specified")
		return
	}

	if host, err := server.db.GetHost(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, host)
	}
}

func (server *EnvServer) handleGetHostsRequest(c *gin.Context) {
	if hosts, err := server.db.GetHosts(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, hosts)
	}
}

func (server *EnvServer) handleRemoveHostRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "Paramater id must be specified")
		return
	}

	if err := server.db.RemoveHost(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, "")
	}
}
