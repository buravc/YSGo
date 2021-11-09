package godis

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

const latestBackupFileName string = "Latest-GodisBackup.json"

type GodisServer struct {
	dictionary *sync.Map
}

func (server *GodisServer) Set(key, val string) {
	server.dictionary.Store(key, val)
}

func (server *GodisServer) Get(key string) (val interface{}, ok bool) {
	return server.dictionary.Load(key)
}

func (server *GodisServer) Flush() {
	server.dictionary.Range(func(key, value interface{}) bool {
		server.dictionary.Delete(key)
		return true
	})
}

func (server *GodisServer) toJson() ([]byte, error) {
	tempMap := make(map[string]string)

	server.dictionary.Range(func(key, value interface{}) bool {
		tempMap[key.(string)] = value.(string)
		return true
	})

	jsonDictionary, err := json.Marshal(tempMap)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return jsonDictionary, nil
}

func (server *GodisServer) fromJson(jsonDictionary []byte) error {
	tempMap := make(map[string]string)

	err := json.Unmarshal(jsonDictionary, &tempMap)

	if err != nil {
		log.Println(err)
		return err
	}

	for k, v := range tempMap {
		server.dictionary.Store(k, v)
	}

	return nil
}

func getLatestBackupFileName() string {
	return os.TempDir() + "/" + latestBackupFileName
}

func (server *GodisServer) backupToFile() (string, error) {

	file, err := os.CreateTemp("", "*-GodisBackup.json")

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	jsonDictionary, err := server.toJson()

	if err != nil {
		os.Remove(file.Name())
		log.Println(err)
	}

	_, err = file.Write(jsonDictionary)

	if err != nil {
		log.Println(err)
		return "", err
	}

	os.Remove(os.TempDir() + "/" + latestBackupFileName)
	err = os.Symlink(file.Name(), getLatestBackupFileName())

	if err != nil {
		return "", nil
	}

	return file.Name(), nil
}

func (server *GodisServer) loadFromFile() {

	server.Flush()

	file, err := os.ReadFile(getLatestBackupFileName())

	if err != nil {
		log.Println(err)
	}

	err = server.fromJson(file)

	if err != nil {
		log.Println(err)
	}

	log.Println("Backup has been restored successfully.")

}
