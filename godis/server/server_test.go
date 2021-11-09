package godis

import (
	"os"
	"testing"
	"time"
)

func Test_backupToFile(t *testing.T) {
	r := buildGodisServer(WithoutBackup, false)
	r.Set("test", "valTest")
	r.Set("test2", "valTest2")
	fileName, err := r.backupToFile()

	if err != nil {
		t.Fail()
	}

	jsnByteArr, err := os.ReadFile(fileName)

	if err != nil {
		t.Fail()
	}

	jsn := string(jsnByteArr)

	jsn2ByteArr, err := r.toJson()

	jsn2 := string(jsn2ByteArr)

	if err != nil {
		t.Fail()
	}

	if jsn != string(jsn2) {
		t.Fail()
	}

	t.Log(jsn)
	t.Log(jsn2)

	t.Run("Test_loadFromFile", Test_loadFromFile)
}

func Test_loadFromFile(t *testing.T) {
	r := buildGodisServer(WithoutBackup, false)

	r.loadFromFile()

	val, ok := r.Get("test")

	if !ok {
		t.Fail()
	}

	if val != "valTest" {
		t.Fail()
	}

	val, ok = r.Get("test2")

	if !ok {
		t.Fail()
	}

	if val != "valTest2" {
		t.Fail()
	}
}

func Test_cronTest(t *testing.T) {
	os.Setenv("GodisBackupInterval", "10s")
	r := buildGodisServer(WithBackup, false)

	r.Set("cronTest", "cronTestVal")

	t.Log("Sleeping 12 seconds")
	time.Sleep(time.Second * 12)

	jsonArr1, err := os.ReadFile(getLatestBackupFileName())

	if err != nil {
		t.Fail()
	}

	jsonStr1 := string(jsonArr1)

	jsonArr2, err := r.toJson()

	if err != nil {
		t.Fail()
	}

	jsonStr2 := string(jsonArr2)

	if jsonStr1 != jsonStr2 {
		t.Fail()
	}
}
