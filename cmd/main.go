package main

import (
	"YSGo/godisApi"
	"os"
)

func main() {
	os.Setenv("GodisBackupInterval", "1m")
	godisApi.NewApiServer()
}
