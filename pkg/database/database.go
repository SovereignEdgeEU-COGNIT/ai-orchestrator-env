package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Postgresql interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type Database struct {
	postgresql  Postgresql
	dbHost      string
	dbPort      int
	dbUser      string
	dbPassword  string
	dbName      string
	dbPrefix    string
	timescaleDB bool
}

func CreateDatabase(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string, dbPrefix string, timescaleDB bool) *Database {
	return &Database{dbHost: dbHost, dbPort: dbPort, dbUser: dbUser, dbPassword: dbPassword, dbName: dbName, dbPrefix: dbPrefix, timescaleDB: timescaleDB}
}

func (db *Database) Connect() error {
	tz := os.Getenv("TZ")
	if tz == "" {
		return errors.New("Timezon environmental variable missing, try e.g. export TZ=Europe/Stockholm")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", db.dbHost, db.dbPort, db.dbUser, db.dbPassword, db.dbName, tz)

	postgresql, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	db.postgresql = postgresql

	err = db.postgresql.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() {
	db.postgresql.Close()
}

func (db *Database) dropHostsTable() error {
	sqlStatement := `DROP TABLE ` + db.dbPrefix + `HOSTS`
	_, err := db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) dropVMsTable() error {
	sqlStatement := `DROP TABLE ` + db.dbPrefix + `VMS`
	_, err := db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Drop() error {
	err := db.dropHostsTable()
	if err != nil {
		return err
	}

	err = db.dropVMsTable()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) createHypertables() error {
	prefix := strings.ToLower(db.dbPrefix)
	sqlStatement := `SELECT create_hypertable ('` + prefix + `hosts', by_range('timestamp', INTERVAL '1 day'))`
	_, err := db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	sqlStatement = `SELECT create_hypertable ('` + prefix + `vms', by_range('timestamp', INTERVAL '1 day'))`
	_, err = db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) createHostsTable() error {
	//sqlStatement := `CREATE TABLE ` + db.dbPrefix + `HOSTS (HOSTID TEXT PRIMARY KEY NOT NULL, HOSTNAME TEXT NOT, TS TIMESTAMPTZ, CPU BIGINT, MEMORY BIGINT)`
	sqlStatement := `CREATE TABLE ` + db.dbPrefix + `HOSTS (HOSTID TEXT PRIMARY KEY NOT NULL, HOSTNAME TEXT NOT NULL)`
	_, err := db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) createVMsTable() error {
	sqlStatement := `CREATE TABLE ` + db.dbPrefix + `VMS (VMID TEXT PRIMARY KEY NOT NULL, VMNAME TEXT NOT NULL)`
	_, err := db.postgresql.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) Initialize() error {
	err := db.createHostsTable()
	if err != nil {
		return err
	}

	err = db.createVMsTable()
	if err != nil {
		return err
	}

	// err = db.createHypertables()
	// if err != nil {
	// 	return err
	// }

	return nil
}
