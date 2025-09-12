//||------------------------------------------------------------------------------------------------||
//|| Log Package: Unit Tests
//|| unit_test.go
//||------------------------------------------------------------------------------------------------||

package log

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

var _ Logger = DefaultLogger{}

//||------------------------------------------------------------------------------------------------||
//|| Logger
//||------------------------------------------------------------------------------------------------||

//||------------------------------------------------------------------------------------------------||
//|| captureStdout runs f while redirecting os.Stdout and returns what was written.
//||------------------------------------------------------------------------------------------------||

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var b strings.Builder
		_, _ = io.Copy(&b, r)
		done <- b.String()
	}()

	// run
	f()

	// restore & close
	_ = w.Close()
	os.Stdout = orig

	select {
	case out := <-done:
		return out
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout reading captured stdout")
		return ""
	}
}

//||------------------------------------------------------------------------------------------------||
//|| captureStderr runs f while redirecting os.Stderr and returns what was written.
//||------------------------------------------------------------------------------------------------||

func captureStderr(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}
	os.Stderr = w

	done := make(chan string, 1)
	go func() {
		var b strings.Builder
		_, _ = io.Copy(&b, r)
		done <- b.String()
	}()

	// run
	f()

	// restore & close
	_ = w.Close()
	os.Stderr = orig

	select {
	case out := <-done:
		return out
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout reading captured stderr")
		return ""
	}
}

// ||------------------------------------------------------------------------------------------------||
// || Return DefautlLogger with Module
// ||------------------------------------------------------------------------------------------------||
func TestInit_ReturnsDefaultLoggerWithModule(t *testing.T) {
	l := Init("api")
	dl, ok := l.(DefaultLogger)
	if !ok {
		t.Fatalf("Init should return DefaultLogger, got %T", l)
	}
	if dl.Module != "api" {
		t.Fatalf("module = %q, want %q", dl.Module, "api")
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test DefaultLogger methods print to correct outputs
//||------------------------------------------------------------------------------------------------||

func TestDefaultLogger_StdoutLevels(t *testing.T) {
	l := DefaultLogger{Module: "testmod"}

	out := captureStdout(t, func() {
		l.Info("hello %s", "world")
		l.Warn("be %s", "careful")
		l.Debug("dbg=%d", 42)
	})

	if !strings.Contains(out, "[INFO]") || !strings.Contains(out, "testmod") || !strings.Contains(out, "hello world") {
		t.Fatalf("INFO not printed correctly: %q", out)
	}
	if !strings.Contains(out, "[WARN]") || !strings.Contains(out, "be careful") {
		t.Fatalf("WARN not printed correctly: %q", out)
	}
	if !strings.Contains(out, "[DEBUG]") || !strings.Contains(out, "dbg=42") {
		t.Fatalf("DEBUG not printed correctly: %q", out)
	}
}

// ||------------------------------------------------------------------------------------------------||
// || Return DefautlLogger Error
// ||------------------------------------------------------------------------------------------------||
func TestDefaultLogger_StderrError(t *testing.T) {
	l := DefaultLogger{Module: "errmod"}

	errOut := captureStderr(t, func() {
		l.Error("boom %d", 99)
	})

	if !strings.Contains(errOut, "[ERROR]") || !strings.Contains(errOut, "errmod") || !strings.Contains(errOut, "boom 99") {
		t.Fatalf("ERROR not printed to stderr correctly: %q", errOut)
	}
}
