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

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `HOSTS (HOSTID, HOSTNAME) VALUES ($1, $2)`
	_, err := db.postgresql.Exec(sqlStatement, host.HostID, host.Hostname)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) parseHosts(rows *sql.Rows) ([]*core.Host, error) {
	var hosts []*core.Host

	for rows.Next() {
		var hostID string
		var hostname string
		if err := rows.Scan(&hostID, &hostname); err != nil {
			return nil, err
		}

		host := &core.Host{HostID: hostID, Hostname: hostname}
		hosts = append(hosts, host)
	}

	return hosts, nil
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
