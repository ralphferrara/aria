package config

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

//||------------------------------------------------------------------------------------------------||
//|| PrintConsoleConfig: Print structured config by section
//||------------------------------------------------------------------------------------------------||

func PrintConsoleConfig(
	c *Config,
) {
	if c == nil {
		fmt.Println("||------------------------------------------------------------------------------------------------||")
		fmt.Println("|| Configuration is nil")
		fmt.Println("||------------------------------------------------------------------------------------------------||")
		return
	}

	fmt.Println("||------------------------------------------------------------------------------------------------||")
	fmt.Println("|| Welcome to Aria")
	fmt.Println("|| Loading Configuration")
	fmt.Println("||------------------------------------------------------------------------------------------------||")
	fmt.Println("")

	val := reflect.ValueOf(*c)
	typ := reflect.TypeOf(*c)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		sectionName := fieldType.Name
		if stringInSlice(sectionName, []string{"App", "Locale", "Template", "Auth"}) {
			fmt.Println("")
			fmt.Printf("\033[1;36m     %s\033[0m\n", sectionName)
			fmt.Println("     -------------------------------------------------------------------------------------------||")
		}
		printJSON(sectionName, fieldVal.Interface())
	}
	fmt.Println("")
	fmt.Println("||------------------------------------------------------------------------------------------------||")
	fmt.Println("|| End Configuration")
	fmt.Println("||------------------------------------------------------------------------------------------------||")
	fmt.Println("")
	fmt.Println("")

}

//||------------------------------------------------------------------------------------------------||
//|| printJSON: Reflectively print any struct or map
//||------------------------------------------------------------------------------------------------||

func printJSON(
	sectionName string,
	data interface{},
) {
	val := reflect.ValueOf(data)

	switch val.Kind() {

	case reflect.Map:
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, key := range keys {
			k := key.String()
			v := val.MapIndex(key)
			fmt.Println()
			fmt.Printf("      %s - \033[1;35m%s\033[0m\n", sectionName, k)
			fmt.Println("      -------------------------------------------------------------------------------------------||")
			printValue("", v.Interface(), 2)
		}

	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := t.Field(i)
			value := val.Field(i)

			if !value.CanInterface() {
				continue
			}

			printValue(field.Name, value.Interface(), 2)
		}

	default:
		fmt.Printf("   %v\n", data)
	}
}

//||------------------------------------------------------------------------------------------------||
//|| printValue: Recursive formatter for all field types
//||------------------------------------------------------------------------------------------------||

func printValue(
	field string,
	data interface{},
	indent int,
) {

	//||------------------------------------------------------------------------------------------------||
	//|| Var
	//||------------------------------------------------------------------------------------------------||

	pad := strings.Repeat("   ", indent)
	val := reflect.ValueOf(data)

	//||------------------------------------------------------------------------------------------------||
	//|| Mask
	//||------------------------------------------------------------------------------------------------||

	if stringInSlice(field, []string{"Password", "AccessKey", "SecretKey"}) {
		fmt.Printf("%s%s : %s\n", pad, field, "*******")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Hide
	//||------------------------------------------------------------------------------------------------||

	if stringInSlice(field, []string{"Middleware", "ErrorHandler", "DB", "Servers", "Dir", "URI", "SSLMode"}) {
		fmt.Printf("%s%s : %s\n", pad, field, "*******")
		return
	}

	//||------------------------------------------------------------------------------------------------||
	//|| Kind
	//||------------------------------------------------------------------------------------------------||

	switch val.Kind() {

	case reflect.Map:
		if field != "" {
			fmt.Printf("%s%s:\n", pad, field)
		}
		keys := val.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, key := range keys {
			fmt.Printf("%s[%s]\n", pad+"   ", key.String())
			printValue("", val.MapIndex(key).Interface(), indent+2)
		}

	case reflect.Struct:
		if field != "" {
			fmt.Printf("%s%s:\n", pad, field)
		}
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			f := t.Field(i)
			v := val.Field(i)
			if v.CanInterface() {
				printValue(f.Name, v.Interface(), indent+1)
			}
		}

	case reflect.Slice:
		if field != "" {
			fmt.Printf("%s%s:\n", pad, field)
		}
		for i := 0; i < val.Len(); i++ {
			printValue("", val.Index(i).Interface(), indent+1)
		}

	default:
		if field != "" {
			fmt.Printf("%s%s : %v\n", pad, field, data)
		} else {
			fmt.Printf("%s%v\n", pad, data)
		}
	}
}

//||------------------------------------------------------------------------------------------------||
//|| StringInSlice: Check if string exists in slice
//||------------------------------------------------------------------------------------------------||

func stringInSlice(
	target string,
	list []string,
) bool {
	for _, item := range list {
		if strings.ToUpper(item) == strings.ToUpper(target) {
			return true
		}
	}
	return false
}
