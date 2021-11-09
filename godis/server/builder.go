package godis

import (
	"log"
	"os"
	"sync"

	"YSGo/godis/cron"
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

		intervalDuration := parseString(intervalStr)

		log.Printf("Backup Interval: %s\n", intervalStr)

		cron.CreateCron(intervalDuration, func() {
			server.backupToFile()
		})
	}

	if loadBackup {
		server.loadFromFile()
	}

	return server

}
