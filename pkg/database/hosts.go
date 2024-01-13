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

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `HOSTS (HOSTID, STATEID, HOSTNAME, CURRENTCPU, CURRENTMEMORY) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.postgresql.Exec(sqlStatement, host.HostID, stateID, host.Hostname, host.CurrentCPU, host.CurrentMemory)
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
		var hostname string
		var currentCPU int64
		var currentMem int64
		if err := rows.Scan(&hostID, &stateID, &hostname, &currentCPU, &currentMem); err != nil {
			return nil, err
		}

		host := &core.Host{HostID: hostID, StateID: stateID, Hostname: hostname, CurrentCPU: currentCPU, CurrentMemory: currentMem}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (db *Database) SetHostResources(hostID string, currentCPU, currentMemory int64) error {
	sqlStatement := `UPDATE ` + db.dbPrefix + `HOSTS SET CURRENTCPU = $1, CURRENTMEMORY = $2 WHERE HOSTID = $3`
	_, err := db.postgresql.Exec(sqlStatement, currentCPU, currentMemory, hostID)
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
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `HOSTS`
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
