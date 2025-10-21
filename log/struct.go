//||------------------------------------------------------------------------------------------------||
//|| Log Package: Structs
//|| struct.go
//||------------------------------------------------------------------------------------------------||

package log

//||------------------------------------------------------------------------------------------------||
//|| Logger Interface (so we can drop in any backend later)
//||------------------------------------------------------------------------------------------------||

type Logger interface {
	Data(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}
