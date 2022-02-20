package jsonlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestLoggerPrintInfo(t *testing.T) {
	err := simulatePrinterLevel(LevelInfo)
	if err != nil {
		t.Error(err)
	}
}

func TestLoggerPrintError(t *testing.T) {
	err := simulatePrinterLevel(LevelError)
	if err != nil {
		t.Error(err)
	}
}

func TestLoggerPrintFatal(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		err := simulatePrinterLevel(LevelFatal)
		if err != nil {
			t.Error(err)
		}
		return
	}

	// We test that a fatal error is raised by invoking go test in a separate process
	cmd := exec.Command(os.Args[0], "-test.run=TestLoggerPrintFatal")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func simulatePrinterLevel(level Level) error {
	tests := []struct {
		message    string
		properties map[string]string
	}{
		{message: "test message"},
		{message: "test message", properties: map[string]string{"env": "test"}},
	}

	for _, tt := range tests {
		b := bytes.NewBuffer(nil)

		var (
			middlewareAfterPrintErrorCalled bool
		)

		logger := New(b, level, Middlewares{
			AfterPrintError: func(err error) {
				middlewareAfterPrintErrorCalled = true
			},
		})

		switch level {
		case LevelInfo:
			logger.PrintInfo(tt.message, tt.properties)
		case LevelError:
			logger.PrintError(errors.New(tt.message), tt.properties)
		case LevelFatal:
			logger.PrintFatal(errors.New(tt.message), tt.properties)
		}

		if level == LevelError && !middlewareAfterPrintErrorCalled {
			return fmt.Errorf("middleware AfterPrintError not called, expected to be called")
		}

		var d details

		err := json.Unmarshal(b.Bytes(), &d)
		if err != nil {
			return err
		}

		if d.Level != level.String() {
			return fmt.Errorf("got %s, expected %s", d.Level, level.String())
		}

		if d.Message != tt.message {
			return fmt.Errorf("got %s, expected %s", d.Message, tt.message)
		}

		if reflect.DeepEqual(d.Properties, tt.properties) != true {
			return fmt.Errorf("got %s, expected %s", d.Properties, tt.properties)
		}
	}

	return nil
}
