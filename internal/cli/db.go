package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbDropCmd)
	rootCmd.AddCommand(dbCmd)

	dbCmd.PersistentFlags().StringVarP(&DBHost, "dbhost", "", DefaultDBHost, "DB host")
	dbCmd.PersistentFlags().IntVarP(&DBPort, "dbport", "", DefaultDBPort, "DB port")
	dbCmd.PersistentFlags().StringVarP(&DBUser, "dbuser", "", "", "DB user")
	dbCmd.PersistentFlags().StringVarP(&DBPassword, "dbpassword", "", "", "DB password")
}

var dbCmd = &cobra.Command{
	Use:   "database",
	Short: "Manage internal database",
	Long:  "Manage internal database",
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database",
	Long:  "Create a database",
	Run: func(cmd *cobra.Command, args []string) {
		parseEnv()
		parseDBEnv()

		var db *database.Database
		for {
			db = database.CreateDatabase(DBHost, DBPort, DBUser, DBPassword, DBName, DBPrefix)
			err := db.Connect()
			if err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Failed to call db.Connect(), retrying in 1 second ...")
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}

		log.WithFields(log.Fields{"Host": DBHost, "Port": DBPort, "User": DBUser, "Password": "**********************", "Prefix": DBPrefix}).Info("Connected to TimescaleDB")

		err := db.Initialize()
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Error("Failed to call db.Initialize()")
			os.Exit(0)
		}

		log.Info("TimescaleDB initialized")
	},
}

var dbDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop the database",
	Long:  "Drop the database",
	Run: func(cmd *cobra.Command, args []string) {
		parseDBEnv()

		fmt.Print("WARNING!!! Are you sure you want to drop the database? This operation cannot be undone! (YES,no): ")

		reader := bufio.NewReader(os.Stdin)
		reply, _ := reader.ReadString('\n')

		if reply == "YES\n" {
			log.WithFields(log.Fields{"DBHost": DBHost, "DBPort": DBPort, "DBUser": DBUser, "DBPassword": "*******************", "DBName": DBName, "UseTLS": UseTLS}).Info("Connecting to TimescaleDB")

			db := database.CreateDatabase(DBHost, DBPort, DBUser, DBPassword, DBName, DBPrefix)
			err := db.Connect()
			CheckError(err)

			err = db.Drop()
			CheckError(err)
			log.Info("TimescaleDB tables dropped")
		} else {
			log.Info("Aborting ...")
		}
	},
}
