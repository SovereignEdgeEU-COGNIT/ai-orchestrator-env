package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
)

func (db *Database) AddMetric(id string, metricType int, metric *core.Metric) error {
	if id == "" {
		return errors.New("ID must be specified")
	}

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `METRICS (HOSTID, TYPE, TS, CPU, MEMORY) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.postgresql.Exec(sqlStatement, id, metricType, metric.Timestamp, metric.CPU, metric.Memory)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) parseMetrics(rows *sql.Rows) ([]*core.Metric, error) {
	var metrics []*core.Metric

	for rows.Next() {
		var hostID string
		var metricType int
		var ts time.Time
		var cpu int64
		var memory int64
		if err := rows.Scan(&hostID, &metricType, &ts, &cpu, &memory); err != nil {
			return nil, err
		}

		metric := &core.Metric{Timestamp: ts, CPU: cpu, Memory: memory}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (db *Database) GetMetrics(hostID string, metricType int, since time.Time, count int) ([]*core.Metric, error) {
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `METRICS WHERE HOSTID=$1 AND TYPE=$2 AND TS>$3 LIMIT $4`
	rows, err := db.postgresql.Query(sqlStatement, hostID, metricType, since, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return db.parseMetrics(rows)
}
