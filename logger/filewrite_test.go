package logger

import (
	"log"
	"os"
	"testing"
)

var filename string

func Testwritetofile(t *testing.T) {
	filename := "testlogfile.txt"
	if !FileExists(filename) {
		CreateFile(filename)
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")
}
