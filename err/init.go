//||------------------------------------------------------------------------------------------------||
//|| Error Package: Registry, Levels, and Helpers
//||------------------------------------------------------------------------------------------------||

package err

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"os"
	"runtime/debug"
)

//||------------------------------------------------------------------------------------------------||
//|| Error Levels
//||------------------------------------------------------------------------------------------------||

const (
	EFATAL   = 0
	EERROR   = 1
	EWARNING = 2
)

//||------------------------------------------------------------------------------------------------||
//|| Error Code Struct
//||------------------------------------------------------------------------------------------------||

type ErrorCode struct {
	Ref     string `json:"ref"`
	Level   int    `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Info    any    `json:"info,omitempty"`
}

//||------------------------------------------------------------------------------------------------||
//|| Error Row Struct
//||------------------------------------------------------------------------------------------------||

type ErrorRow struct {
	Code    string
	Message string
	Level   int
}

//||------------------------------------------------------------------------------------------------||
//|| Error Library Struct
//||------------------------------------------------------------------------------------------------||

type ErrorLibrary struct {
	Reference string
	Entries   map[string]ErrorCode
}

//||------------------------------------------------------------------------------------------------||
//|| Error Library Registry (by module/category)
//||------------------------------------------------------------------------------------------------||

var (
	ERR           = map[string]*ErrorLibrary{}
	CurrentModule string
)

//||------------------------------------------------------------------------------------------------||
//|| GetModuleName: Returns the go.mod module name at runtime
//||------------------------------------------------------------------------------------------------||

func GetModuleName() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Path
	}
	return "UNKNOWN"
}

//||------------------------------------------------------------------------------------------------||
//|| AddLibrary: Register a new error library/category (returns pointer)
//||------------------------------------------------------------------------------------------------||

func AddLibrary(ref string) *ErrorLibrary {
	if lib, exists := ERR[ref]; exists {
		return lib
	}
	lib := &ErrorLibrary{
		Reference: ref,
		Entries:   map[string]ErrorCode{},
	}
	ERR[ref] = lib
	return lib
}

//||------------------------------------------------------------------------------------------------||
//|| RegisterError: Register a single error to current module library
//||------------------------------------------------------------------------------------------------||

func RegisterError(code, message string, level int) {
	lib := AddLibrary(CurrentModule)
	lib.RegisterError(code, message, level)
}

//||------------------------------------------------------------------------------------------------||
//|| RegisterErrors: Register multiple errors to current module library
//||------------------------------------------------------------------------------------------------||

func RegisterErrors(rows []ErrorRow) {
	lib := AddLibrary(CurrentModule)
	for _, row := range rows {
		lib.RegisterError(row.Code, row.Message, row.Level)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| RegisterError (instance): Add to Entries
//||------------------------------------------------------------------------------------------------||

func (lib *ErrorLibrary) RegisterError(code, message string, level int) {
	lib.Entries[code] = ErrorCode{
		Ref:     lib.Reference,
		Code:    code,
		Message: message,
		Level:   level,
	}
}

//||------------------------------------------------------------------------------------------------||
//|| Debug Method
//||------------------------------------------------------------------------------------------------||

func (e ErrorCode) Debug(data any) ErrorCode {
	e.Info = data
	if os.Getenv("ENV_MODE") == "development" {
		if e.Level == EFATAL {
			fmt.Println("Stack Trace:", string(debug.Stack()))
			fmt.Println("Error Info:", e.Info)
			panic("Fatal error: " + e.Message)
		} else {
			fmt.Printf("Error Code: %s, Message: %s, Level: %d, Info: %v\n", e.Code, e.Message, e.Level, e.Info)
		}
		return e
	}
	// TO DO: Production logging
	return e
}

//||------------------------------------------------------------------------------------------------||
//|| Error interface: Returns a standard error
//||------------------------------------------------------------------------------------------------||

func (e ErrorCode) Error() error {
	return fmt.Errorf("%s.%s:%s", e.Ref, e.Code, e.Message)
}

//||------------------------------------------------------------------------------------------------||
//|| Panic: Panics with the error
//||------------------------------------------------------------------------------------------------||

func (e ErrorCode) Panic() error {
	panic(e.Error())
}
