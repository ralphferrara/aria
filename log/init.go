//||------------------------------------------------------------------------------------------------||
//|| Log Package: Initialization (All-in-One)
//|| init.go
//||------------------------------------------------------------------------------------------------||

package log

import (
	"fmt"
	"os"
	"time"

	"aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Log: Level Type and Constants
//||------------------------------------------------------------------------------------------------||

type Level string

const (
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	DEBUG Level = "DEBUG"
)

//||------------------------------------------------------------------------------------------------||
//|| Log: Global Config Reference
//||------------------------------------------------------------------------------------------------||

var AppConfig *config.Config

func Init(cfg *config.Config) {
	AppConfig = cfg
}

//||------------------------------------------------------------------------------------------------||
//|| Log: Print (General Logger)
//||------------------------------------------------------------------------------------------------||

func Print(level Level, module string, msg string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fullMsg := fmt.Sprintf(msg, args...)
	line := fmt.Sprintf("[%s] [%s] [%s] %s\n", now, level, module, fullMsg)
	if level == ERROR {
		fmt.Fprint(os.Stderr, line)
	} else {
		fmt.Print(line)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Log: Level-Specific Shortcuts
//||------------------------------------------------------------------------------------------------||

func Info(module, msg string, args ...interface{})  { Print(INFO, module, msg, args...) }
func Warn(module, msg string, args ...interface{})  { Print(WARN, module, msg, args...) }
func Error(module, msg string, args ...interface{}) { Print(ERROR, module, msg, args...) }
func Debug(module, msg string, args ...interface{}) { Print(DEBUG, module, msg, args...) }

//||------------------------------------------------------------------------------------------------||
//|| Log: Facade Struct for Log.Info-style calls
//||------------------------------------------------------------------------------------------------||

var Log = struct {
	Info      func(module, msg string, args ...interface{})
	Warn      func(module, msg string, args ...interface{})
	Error     func(module, msg string, args ...interface{})
	Debug     func(module, msg string, args ...interface{})
	Init      func(cfg *config.Config)
	AppConfig **config.Config
}{
	Info:      Info,
	Warn:      Warn,
	Error:     Error,
	Debug:     Debug,
	Init:      Init,
	AppConfig: &AppConfig,
}
