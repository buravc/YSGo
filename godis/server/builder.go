package godis

import (
	"log"
	"os"
	"sync"

	"github.com/callduckk/YSGo/godis/cron"
)

type ServerType int

const (
	WithBackup ServerType = iota + 1
	WithoutBackup
)

func buildGodisServer(serverType ServerType, loadBackup bool) *GodisServer {

	server := &GodisServer{}
	server.dictionary = &sync.Map{}

	if serverType == WithBackup {
		intervalStr := os.Getenv("GodisBackupInterval")

		if intervalStr == "" {
			log.Println("'GodisBackupInterval' env variable was not set. Falling back to 30m")
			intervalStr = "30m"
		}

		intervalDuration := parseString(intervalStr)

		log.Printf("Backup Interval: %s\n", intervalStr)

		log.Printf("Backup path: %s\n", getLatestBackupFileName())

		cron.CreateCron(intervalDuration, func() {
			log.Println("Backing up to file.")
			server.backupToFile()
		})
		log.Println("Backup cron has been created successfully.")
	}

	if loadBackup {
		log.Println("Restoring latest backup.")
		server.loadFromFile()
	}

	return server

}
