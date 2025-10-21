//||------------------------------------------------------------------------------------------------||
//|| Log Package: Simple Wrapper (fmt.Println based)
//|| init.go
//||------------------------------------------------------------------------------------------------||

package log

import (
	"fmt"
	"os"
)

//||------------------------------------------------------------------------------------------------||
//|| Log Message
//||------------------------------------------------------------------------------------------------||

type DefaultLogger struct {
	Module string
}

const BoldYellow = "\033[1;33m"
const Reset = "\033[0m"

//||------------------------------------------------------------------------------------------------||
//|| Defined
//||------------------------------------------------------------------------------------------------||

func (d DefaultLogger) Data(msg string, args ...interface{}) {
	fmt.Println(BoldYellow + fmt.Sprintf("[DATA] %s %s", d.Module, fmt.Sprintf(msg, args...)) + Reset)
}

func (d DefaultLogger) Info(msg string, args ...interface{}) {
	fmt.Println("[INFO]", d.Module, fmt.Sprintf(msg, args...))
}
func (d DefaultLogger) Warn(msg string, args ...interface{}) {
	fmt.Println("[WARN]", d.Module, fmt.Sprintf(msg, args...))
}
func (d DefaultLogger) Error(msg string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, "[ERROR]", d.Module, fmt.Sprintf(msg, args...))
}
func (d DefaultLogger) Debug(msg string, args ...interface{}) {
	fmt.Println("[DEBUG]", d.Module, fmt.Sprintf(msg, args...))
}
