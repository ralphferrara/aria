package log

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/ralphferrara/aria/config"
)

func TestInitSetsAppConfig(t *testing.T) {
	cfg := &config.Config{}
	Init(cfg)
	if AppConfig != cfg {
		t.Error("AppConfig was not set by Init")
	}
}

// Helper to capture output
func captureOutput(f func()) (string, string) {
	// Save current stdout/stderr
	oldOut, oldErr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout, os.Stderr = wOut, wErr

	// Call function
	f()

	// Restore and read
	wOut.Close()
	wErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var bufOut, bufErr bytes.Buffer
	bufOut.ReadFrom(rOut)
	bufErr.ReadFrom(rErr)
	return bufOut.String(), bufErr.String()
}

func TestPrintAndShortcuts(t *testing.T) {
	// Test INFO output to stdout
	out, err := captureOutput(func() {
		Info("testmod", "Info message: %d", 42)
		Warn("testmod", "Warn message: %s", "be careful")
		Debug("testmod", "Debug message!")
	})
	if !strings.Contains(out, "[INFO] [testmod] Info message: 42") {
		t.Error("Info() output missing or malformed")
	}
	if !strings.Contains(out, "[WARN] [testmod] Warn message: be careful") {
		t.Error("Warn() output missing or malformed")
	}
	if !strings.Contains(out, "[DEBUG] [testmod] Debug message!") {
		t.Error("Debug() output missing or malformed")
	}
	if err != "" {
		t.Errorf("Did not expect output to stderr for INFO/WARN/DEBUG, got: %q", err)
	}

	// Test ERROR output to stderr
	out, err = captureOutput(func() {
		Error("errmod", "Error message: %s", "fail!")
	})
	if !strings.Contains(err, "[ERROR] [errmod] Error message: fail!") {
		t.Error("Error() output missing or malformed")
	}
	if out != "" {
		t.Errorf("Did not expect output to stdout for ERROR, got: %q", out)
	}
}

func TestFacadeStruct(t *testing.T) {
	out, _ := captureOutput(func() {
		Log.Info("facade", "test %s", "one")
		Log.Warn("facade", "test %d", 2)
		Log.Debug("facade", "test")
	})
	if !strings.Contains(out, "test one") {
		t.Error("Log.Info did not output expected message")
	}
}
