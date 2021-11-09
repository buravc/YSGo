package godis

var instance *godisServer

func Singleton() *godisServer {

	if instance == nil {
		instance = buildGodisServer(WithBackup, true)
	}

	return instance
}
