package server

import (
	"io"
	"net/http"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/gin-gonic/gin"
)

func (server *EnvServer) handleAddVMRequest(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	vm, err := core.ConvertJSONToVM(string(jsonData))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = server.db.AddVM(vm)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "")
}

func (server *EnvServer) handleGetVMRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "Paramater id must be specified")
		return
	}

	if vm, err := server.db.GetVM(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, vm)
	}
}

func (server *EnvServer) handleGetVMsRequest(c *gin.Context) {
	if vms, err := server.db.GetVMs(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, vms)
	}
}

func (server *EnvServer) handleRemoveVMRequest(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "Paramater id must be specified")
		return
	}

	if err := server.db.RemoveVM(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, "")
	}
}
