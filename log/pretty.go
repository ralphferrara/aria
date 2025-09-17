package log

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "   ")
	if err != nil {
		fmt.Println("Error pretty printing:", err)
		return
	}
	fmt.Println(string(bytes))
}
