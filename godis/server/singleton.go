package godis

import (
	"sync"
)

var instance *GodisServer
var mu sync.Mutex

// we need thread safety here to prevent data race which could result in duplicate godisServer instances if occurred
// a simple lock mechanism should work therefore a mutex is enough.

func Singleton() *GodisServer {

	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = buildGodisServer(WithBackup, true)
	}

	return instance
}
