//||------------------------------------------------------------------------------------------------||
//|| Log Package: Log Functions
//|| load.go
//||------------------------------------------------------------------------------------------------||

package log

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"os"
	"time"
)

//||------------------------------------------------------------------------------------------------||
//|| Log: Print (General Logger)
//||------------------------------------------------------------------------------------------------||

func Print(level Level, module string, msg string, args ...interface{}) {
	if !shouldLog(level) {
		return
	}
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
//|| Log: Level Control Helper
//||------------------------------------------------------------------------------------------------||

func shouldLog(level Level) bool {
	// Simple example: log all for now (expand for filtering if needed)
	return true
}
