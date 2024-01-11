package database

import (
	"database/sql"
	"errors"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/core"
)

func (db *Database) AddVM(host *core.VM) error {
	if host == nil {
		return errors.New("VM is nil")
	}

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `VMS (VMID, HOSTNAME) VALUES ($1, $2)`
	_, err := db.postgresql.Exec(sqlStatement, host.VMID, host.Hostname)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) parseVMs(rows *sql.Rows) ([]*core.VM, error) {
	var vms []*core.VM

	for rows.Next() {
		var vmID string
		var hostname string
		if err := rows.Scan(&vmID, &hostname); err != nil {
			return nil, err
		}

		vm := &core.VM{VMID: vmID, Hostname: hostname}
		vms = append(vms, vm)
	}

	return vms, nil
}

func (db *Database) GetVMs() ([]*core.VM, error) {
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `VMS`
	rows, err := db.postgresql.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return db.parseVMs(rows)
}
