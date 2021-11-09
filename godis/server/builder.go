package godis

import (
	"log"
	"os"
	"sync"
	"time"

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
		var intervalDuration time.Duration
		switch intervalStr {
		case "1h":
			intervalDuration = time.Hour
		case "15m":
			intervalDuration = time.Minute * 15
		case "10s":
			intervalDuration = time.Second * 10
		default:
			intervalDuration = time.Hour
		}

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
