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

	vms, err := db.GetVMs()
	if err != nil {
		return err
	}

	// Create a map to track used stateIDs
	usedStateIDs := make(map[int]bool)
	for _, h := range vms {
		usedStateIDs[h.StateID] = true
	}

	// Find the first available stateID
	stateID := 0
	for usedStateIDs[stateID] {
		stateID++
	}

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `VMS (VMID, STATEID, HOSTNAME, CURRENTCPU, CURRENTMEMORY) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.postgresql.Exec(sqlStatement, host.VMID, stateID, host.Hostname, host.CurrentCPU, host.CurrentMemory)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) parseVMs(rows *sql.Rows) ([]*core.VM, error) {
	var vms []*core.VM

	for rows.Next() {
		var vmID string
		var stateID int
		var hostname string
		var currentCPU int64
		var currentMemory int64
		if err := rows.Scan(&vmID, &stateID, &hostname, &currentCPU, &currentMemory); err != nil {
			return nil, err
		}

		vm := &core.VM{VMID: vmID, StateID: stateID, Hostname: hostname, CurrentCPU: currentCPU, CurrentMemory: currentMemory}
		vms = append(vms, vm)
	}

	return vms, nil
}

func (db *Database) SetVMResources(vmID string, currentCPU, currentMemory int64) error {
	sqlStatement := `UPDATE ` + db.dbPrefix + `VMS SET CURRENTCPU = $1, CURRENTMEMORY = $2 WHERE VMID = $3`
	_, err := db.postgresql.Exec(sqlStatement, currentCPU, currentMemory, vmID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetVM(vmID string) (*core.VM, error) {
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `VMS WHERE VMID = $1`
	rows, err := db.postgresql.Query(sqlStatement, vmID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	vms, err := db.parseVMs(rows)
	if err != nil {
		return nil, err
	}

	if len(vms) == 0 {
		return nil, nil
	}

	return vms[0], nil
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

func (db *Database) RemoveVM(vmID string) error {
	sqlStatement := `DELETE FROM ` + db.dbPrefix + `VMS WHERE VMID=$1`
	_, err := db.postgresql.Exec(sqlStatement, vmID)
	if err != nil {
		return err
	}

	return nil
}
