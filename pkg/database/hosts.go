package database

import (
	"database/sql"
	"errors"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
)

func (db *Database) AddHost(host *core.Host) error {
	if host == nil {
		return errors.New("Host is nil")
	}

	db.hostsMutex.Lock()

	hosts, err := db.GetHosts()
	if err != nil {
		return err
	}

	// Create a map to track used stateIDs
	usedStateIDs := make(map[int]bool)
	for _, h := range hosts {
		usedStateIDs[h.StateID] = true
	}

	// Find the first available stateID, starting at 1
	stateID := 1
	for usedStateIDs[stateID] {
		stateID++
	}

	db.hostsMutex.Unlock()

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `HOSTS (HOSTID, STATEID, TOTAL_CPU, TOTAL_MEM, USAGE_CPU, USAGE_MEM) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.postgresql.Exec(sqlStatement, host.HostID, stateID, host.TotalCPU, host.TotalMemory, host.UsageCPU, host.UsageMemory)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) parseHosts(rows *sql.Rows) ([]*core.Host, error) {
	var hosts []*core.Host

	for rows.Next() {
		var hostID string
		var stateID int
		var totalCPU int64
		var totalMem int64
		var usageCPU int64
		var usageMem int64

		if err := rows.Scan(&hostID, &stateID, &totalCPU, &totalMem, &usageCPU, &usageMem); err != nil {
			return nil, err
		}

		host := &core.Host{
			HostID:      hostID,
			StateID:     stateID,
			TotalCPU:    totalCPU,
			TotalMemory: totalMem,
			UsageCPU:    usageCPU,
			UsageMemory: usageMem,
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (db *Database) SetHostResources(hostID string, usageCPU, usageMemory int64) error {
	sqlStatement := `UPDATE ` + db.dbPrefix + `HOSTS SET USAGE_CPU = $1, USAGE_MEM = $2 WHERE HOSTID = $3`
	_, err := db.postgresql.Exec(sqlStatement, usageCPU, usageMemory, hostID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetHost(hostID string) (*core.Host, error) {
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `HOSTS WHERE HOSTID = $1`
	rows, err := db.postgresql.Query(sqlStatement, hostID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	hosts, err := db.parseHosts(rows)
	if err != nil {
		return nil, err
	}

	if len(hosts) == 0 {
		return nil, nil
	}

	return hosts[0], nil
}

func (db *Database) GetHosts() ([]*core.Host, error) {
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `HOSTS ORDER BY STATEID ASC`
	rows, err := db.postgresql.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return db.parseHosts(rows)
}

func (db *Database) RemoveHost(hostID string) error {
	sqlStatement := `DELETE FROM ` + db.dbPrefix + `HOSTS WHERE HOSTID=$1`
	_, err := db.postgresql.Exec(sqlStatement, hostID)
	if err != nil {
		return err
	}

	return nil
}
