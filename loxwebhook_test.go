package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_initLogging(t *testing.T) {
	logFormat := log.Ldate | log.Ltime | log.Lshortfile
	logFileName := "./test_initLogging.log"
	randomID := rand.Intn(10000)
	logFileContent := fmt.Sprintf("Test_initLogging %d", randomID)

	// Test stderr
	// cfg := &config.Config{
	// 	LogFile: "",
	// }
	t.Run("stderr", func(t *testing.T) {
		wantLogger := log.New(
			os.Stderr,
			"",
			logFormat,
		)
		gotLogger, gotFile, err := initLogging("")
		if err != nil {
			t.Errorf("initLogging() error = %v", err)
			return
		}
		if !reflect.DeepEqual(gotLogger, wantLogger) {
			t.Errorf("initLogging() got = %v, want %v", gotLogger, wantLogger)
		}
		if gotFile != nil {
			t.Errorf("initLogging() got1 = %v, want %v", gotFile, nil)
		}
	})

	// Test file
	// cfg = &config.Config{
	// 	LogFile: logFileName,
	// }
	t.Run("file", func(t *testing.T) {
		gotLogger, gotFile, err := initLogging(logFileName)
		gotLogger.Print(logFileContent)
		defer func() {
			gotLogger.SetFlags(0)
			gotLogger.SetOutput(ioutil.Discard)
			gotFile.Close()
			err = os.Remove(logFileName)
			if err != nil {
				t.Errorf("Deleting logfile failed: %v", err)
			}
		}()
		if err != nil {
			t.Errorf("initLogging() error = %v", err)
			return
		}
		if gotFile.Name() != logFileName {
			t.Errorf("Wrong filename. Expected: %v, got: %v", logFileName, gotFile.Name())
		}
		content, err := ioutil.ReadFile(logFileName)
		if err != nil {
			t.Errorf("Cannot read from logfile")
		}
		if !strings.Contains(string(content), logFileContent) {
			t.Errorf("Wrong content in logfile. Expected: %v, got: %v", logFileContent, string(content))
		}
	})
}
