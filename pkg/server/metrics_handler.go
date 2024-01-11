package server

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
	"github.com/gin-gonic/gin"
)

func (server *EnvServer) handleAddMetricRequest(c *gin.Context) {
	hostID, ok := c.GetQuery("hostid")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater hostid must be specified")
		return
	}

	metricTypeStr, ok := c.GetQuery("metrictype")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater since must be specified")
		return
	}
	metricType, err := strconv.Atoi(metricTypeStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Paramater metrictype must be an integer")
		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	metric, err := core.ConvertJSONToMetric(string(jsonData))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = server.db.AddMetric(hostID, metricType, metric)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, "")
}

func (server *EnvServer) handleGetMetricsRequest(c *gin.Context) {
	hostID, ok := c.GetQuery("hostid")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater hostid must be specified")
		return
	}

	metricTypeStr, ok := c.GetQuery("metrictype")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater since must be specified")
		return
	}
	metricType, err := strconv.Atoi(metricTypeStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Paramater metrictype must be an integer")
		return
	}

	nanoUnixTimeStr, ok := c.GetQuery("since")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater hostid must be specified")
		return
	}

	nanoUnixTime, err := strconv.ParseInt(nanoUnixTimeStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Paramater since must be an integer")
		return
	}

	seconds := nanoUnixTime / int64(time.Second)
	nanoseconds := nanoUnixTime % int64(time.Second)
	ts := time.Unix(seconds, nanoseconds)

	countStr, ok := c.GetQuery("count")
	if !ok {
		c.String(http.StatusBadRequest, "Paramater count must be specified")
		return
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Paramater count must be an integer")
		return
	}

	metrics, err := server.db.GetMetrics(hostID, metricType, ts, count)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	metricsJSON, err := core.ConvertMetricArrayToJSON(metrics)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, metricsJSON)
}
