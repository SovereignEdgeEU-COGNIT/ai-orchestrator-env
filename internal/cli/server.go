package cli

import (
	"time"

	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/database"
	"github.com/SovereignEdgeEU-COGNIT/ai-orchestrator-env/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage Env server",
	Long:  "Manage Env server",
}

func init() {
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().BoolVarP(&InitDB, "initdb", "", false, "Initialize DB")
	serverCmd.PersistentFlags().IntVarP(&ServerPort, "port", "", -1, "Server HTTP port")
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Env server",
	Long:  "Start a Env server",
	Run: func(cmd *cobra.Command, args []string) {
		parseDBEnv()
		parseEnv()

		var db *database.Database
		for {
			db = database.CreateDatabase(DBHost, DBPort, DBUser, DBPassword, DBName, DBPrefix)
			err := db.Connect()
			if err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Failed to connect to PostgreSQL database")
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}

		log.WithFields(log.Fields{"DBHost": DBHost, "DBPort": DBPort, "DBUser": DBUser, "DBPassword": "*******************", "DBName": DBName}).Info("Connected to PostgreSQL database")

		server := server.CreateEnvServer(db, ServerPort)

		if InitDB {
			err := db.Initialize()
			CheckError(err)
		}

		for {
			err := server.ServeForever()
			if err != nil {
				log.WithFields(log.Fields{"Error": err}).Error("Failed to start Env server")
				time.Sleep(1 * time.Second)
			}
		}
	},
}
