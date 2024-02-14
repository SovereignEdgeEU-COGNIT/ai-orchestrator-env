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

	db.vmsMutex.Lock()

	vms, err := db.GetVMs()
	if err != nil {
		return err
	}

	// Create a map to track used stateIDs
	usedStateIDs := make(map[int]bool)
	for _, h := range vms {
		usedStateIDs[h.StateID] = true
	}

	// Find the first available stateID, starting at 1
	stateID := 1
	for usedStateIDs[stateID] {
		stateID++
	}

	db.vmsMutex.Unlock()

	sqlStatement := `INSERT INTO ` + db.dbPrefix + `VMS (VMID, STATEID, DEPLOYED, HOSTID, HOSTSTATEID, TOTAL_CPU, TOTAL_MEM, USAGE_CPU, USAGE_MEM) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.postgresql.Exec(sqlStatement, host.VMID, stateID, false, "", 0, host.TotalCPU, host.TotalMemory, host.UsageCPU, host.UsageMemory)
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
		var deployed bool
		var hostID string
		var hostStateID int
		var totalCPU float64
		var totalMemory int64
		var usageCPU float64
		var usageMemory int64
		if err := rows.Scan(&vmID, &stateID, &deployed, &hostID, &hostStateID, &totalCPU, &totalMemory, &usageCPU, &usageMemory); err != nil {
			return nil, err
		}

		vm := &core.VM{VMID: vmID, StateID: stateID, Deployed: deployed, HostID: hostID, HostStateID: hostStateID, TotalCPU: totalCPU, TotalMemory: totalMemory, UsageCPU: usageCPU, UsageMemory: usageMemory}

		vms = append(vms, vm)
	}

	return vms, nil
}

func (db *Database) SetVMResources(vmID string, usageCPU float64, usageMemory int64) error {
	sqlStatement := `UPDATE ` + db.dbPrefix + `VMS SET USAGE_CPU = $1, USAGE_MEM = $2 WHERE VMID = $3`
	_, err := db.postgresql.Exec(sqlStatement, usageCPU, usageMemory, vmID)
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
	sqlStatement := `SELECT * FROM ` + db.dbPrefix + `VMS ORDER BY STATEID ASC`
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

func (db *Database) Bind(vmID, hostID string) error {
	host, err := db.GetHost(hostID)
	if err != nil {
		return err
	}

	if host == nil {
		return errors.New("Host not found")
	}

	sqlStatement := `UPDATE ` + db.dbPrefix + `VMS SET DEPLOYED = $1, HOSTID = $2, HOSTSTATEID = $3 WHERE VMID = $4`
	_, err = db.postgresql.Exec(sqlStatement, true, host.HostID, host.StateID, vmID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Unbind(vmID string) error {
	sqlStatement := `UPDATE ` + db.dbPrefix + `VMS SET DEPLOYED = $1, HOSTID = $2, HOSTSTATEID = $3 WHERE VMID = $4`
	_, err := db.postgresql.Exec(sqlStatement, false, "", 0, vmID)
	if err != nil {
		return err
	}

	return nil
}
