package database

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
)

func (db *Database) AddMetric(id string, metricType int, metric *core.Metric) error {
	if id == "" {
		return errors.New("ID must be specified")
	}

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `METRICS (ID, TYPE, TS, CPU, MEMORY) VALUES ($1, $2, $3, $4, $5)`
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
		var cpu float64
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
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `METRICS WHERE ID=$1 AND TYPE=$2 AND TS>$3 LIMIT $4`
	rows, err := db.postgresql.Query(sqlStatement, hostID, metricType, since, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return db.parseMetrics(rows)
}

func (db *Database) Export(id string, metricType int, filename string) error {
	sqlStatement := `SELECT ID, TYPE, TS, CPU, MEMORY FROM ` + db.dbPrefix + `METRICS WHERE ID=$1 AND TYPE=$2`
	rows, err := db.postgresql.Query(sqlStatement, id, metricType)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"id", "type", "ts", "cpu", "mem"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write rows to CSV
	for rows.Next() {
		var id string
		var ts time.Time
		var typeValue int
		var cpu float64
		var memory int64

		err := rows.Scan(&id, &typeValue, &ts, &cpu, &memory)
		if err != nil {
			return err
		}

		record := []string{id, fmt.Sprintf("%d", typeValue), strconv.FormatInt(ts.UnixNano(), 10), fmt.Sprintf("%f", cpu), fmt.Sprintf("%d", memory)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
