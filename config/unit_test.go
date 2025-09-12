//||------------------------------------------------------------------------------------------------||
//|| Config Package
//|| unit_test.go
//||------------------------------------------------------------------------------------------------||

package config

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

//||------------------------------------------------------------------------------------------------||
//|| helper: write JSON to a temp file
//||------------------------------------------------------------------------------------------------||

func writeJSON(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

//||------------------------------------------------------------------------------------------------||
//|| Test Init and GetConfig
//||------------------------------------------------------------------------------------------------||

func TestInitAndGetConfig_Basic(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	json := `{
	  "app": { "name":"aria","env":"development","debug":true,"port":8080 },
	  "db": { "main": { "driver":"postgres","host":"localhost","port":5432,"user":"u","password":"p","database":"d","sslmode":"disable" } },
	  "cache": { "primary": { "backend":"redis","host":"localhost","port":6379,"db":0 } },
	  "storage": { "uploads": { "backend":"local","dir":"./uploads" } },
	  "queue": { "main": { "backend":"rabbitmq","host":"localhost","port":5672,"user":"u","password":"p","vhost":"/" } },
	  "http": {
	    "api":    { "backend":"mux","port":8081,"cors":true,"middleware":true,"error_handler":true },
	    "client": { "backend":"HTTP","port":8082,"cors":false,"middleware":false,"error_handler":false }
	  },
	  "auth": { "jwt_secret":"s","session_expiry":3600,"token_issuer":"aria" },
	  "locale": { "default":"en-US","supported":["en-US"] },
	  "template": { "dir":"./templates","cache":true }
	}`
	path := writeJSON(t, tmp, "config.json", json)

	cfg, err := Init(path)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if cfg == nil {
		t.Fatal("Init returned nil")
	}
	if GetConfig() != cfg {
		t.Fatal("GetConfig did not return same instance")
	}
	// normalized http backend
	if got := strings.ToLower(cfg.HTTP["client"].Backend); got != "http" {
		t.Fatalf("normalize(http.backend) = %q, want http", got)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test Init Disallow
//||------------------------------------------------------------------------------------------------||

func TestInit_DisallowUnknownFields(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	// "oops" is unknown
	json := `{
	  "app":{"name":"aria","env":"dev","debug":false,"port":8080},
	  "db":{}, "cache":{}, "storage":{}, "queue":{}, "http":{},
	  "auth":{"jwt_secret":"x","session_expiry":1,"token_issuer":"x"},
	  "locale":{"default":"en-US","supported":[]},
	  "template":{"dir":".","cache":false},
	  "oops": 123
	}`
	path := writeJSON(t, tmp, "bad.json", json)
	if _, err := Init(path); err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test Init Env Expansion
//||------------------------------------------------------------------------------------------------||

func TestInit_EnvExpansion(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	_ = os.Setenv("USER", "aria-test")
	_ = os.Setenv("GOOS", runtime.GOOS) // ensures ${GOOS} expands
	_ = os.Setenv("TEST_SECRET", "abc123")
	_ = os.Setenv("TEST_DIR", "tpls")
	json := `{
	  "app":{"name":"${USER}","env":"${GOOS}","debug":false,"port":8080},
	  "db":{"main":{"driver":"postgres","host":"localhost","port":5432,"user":"u","password":"${TEST_SECRET}","database":"d"}},
	  "cache":{}, "storage":{}, "queue":{}, "http":{},
	  "auth":{"jwt_secret":"${TEST_SECRET}","session_expiry":1,"token_issuer":"aria"},
	  "locale":{"default":"en-US","supported":[]},
	  "template":{"dir":"./${TEST_DIR}","cache":false}
	}`
	path := writeJSON(t, tmp, "env.json", json)
	cfg, err := Init(path)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if cfg.DB["main"].Password != "abc123" {
		t.Fatalf("env expand failed, got %q", cfg.DB["main"].Password)
	}
	if cfg.Template.Dir != "./tpls" {
		t.Fatalf("template dir expand failed, got %q", cfg.Template.Dir)
	}
	// sanity: app fields expanded (best-effort)
	if cfg.App.Env != runtime.GOOS {
		t.Fatalf("app.env expand failed, got %q want %q", cfg.App.Env, runtime.GOOS)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test Reset
//||------------------------------------------------------------------------------------------------||

func TestReset_AllowsReinit(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	json1 := `{
	  "app":{"name":"a","env":"d","debug":false,"port":1001},
	  "db":{}, "cache":{}, "storage":{}, "queue":{}, "http":{},
	  "auth":{"jwt_secret":"s","session_expiry":1,"token_issuer":"a"},
	  "locale":{"default":"en-US","supported":[]},
	  "template":{"dir":".","cache":false}
	}`
	p1 := writeJSON(t, tmp, "c1.json", json1)
	cfg1, err := Init(p1)
	if err != nil {
		t.Fatalf("Init1: %v", err)
	}
	if cfg1.App.Port != 1001 {
		t.Fatalf("port1 = %d", cfg1.App.Port)
	}

	Reset()

	json2 := strings.ReplaceAll(json1, "1001", "2002")
	p2 := writeJSON(t, tmp, "c2.json", json2)
	cfg2, err := Init(p2)
	if err != nil {
		t.Fatalf("Init2: %v", err)
	}
	if cfg2.App.Port != 2002 {
		t.Fatalf("port2 = %d", cfg2.App.Port)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test Panic Init
//||------------------------------------------------------------------------------------------------||

func TestMust_PanicsWithoutInit(t *testing.T) {
	t.Cleanup(Reset)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Must() should panic without Init")
		}
	}()
	_ = Must()
}

//||------------------------------------------------------------------------------------------------||
//|| Port Validation
//||------------------------------------------------------------------------------------------------||

func TestHTTPPortValidation(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	json := `{
	  "app":{"name":"a","env":"d","debug":false,"port":8080},
	  "db":{}, "cache":{}, "storage":{}, "queue":{},
	  "http":{"api":{"backend":"mux","port":70000,"cors":false,"middleware":false,"error_handler":false}},
	  "auth":{"jwt_secret":"s","session_expiry":1,"token_issuer":"a"},
	  "locale":{"default":"en-US","supported":[]},
	  "template":{"dir":".","cache":false}
	}`
	path := writeJSON(t, tmp, "badport.json", json)
	if _, err := Init(path); err == nil {
		t.Fatal("expected error for invalid http port, got nil")
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Test DB
//||------------------------------------------------------------------------------------------------||

func TestDBUnsupportedDriver(t *testing.T) {
	t.Cleanup(Reset)

	tmp := t.TempDir()
	json := `{
	  "app":{"name":"a","env":"d","debug":false,"port":8080},
	  "db":{"x":{"driver":"oracle"}},"cache":{}, "storage":{}, "queue":{}, "http":{},
	  "auth":{"jwt_secret":"s","session_expiry":1,"token_issuer":"a"},
	  "locale":{"default":"en-US","supported":[]},
	  "template":{"dir":".","cache":false}
	}`
	path := writeJSON(t, tmp, "baddriver.json", json)
	if _, err := Init(path); err == nil {
		t.Fatal("expected error for unsupported db driver, got nil")
	}
}
