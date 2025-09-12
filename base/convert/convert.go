package convert

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

//||------------------------------------------------------------------------------------------------||
//|| Type IntBool
//||------------------------------------------------------------------------------------------------||

type IntBool int

func (ib *IntBool) UnmarshalJSON(data []byte) error {
	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s == "1" {
			*ib = 1
		} else {
			*ib = 0
		}
		return nil
	}

	// Try number
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*ib = IntBool(i)
		return nil
	}

	return fmt.Errorf("invalid IntBool input: %s", string(data))
}

//||------------------------------------------------------------------------------------------------||
//|| Convert int8 to pointer
//||------------------------------------------------------------------------------------------------||

func Int8Ptr(v int8) *int8 {
	return &v
}

//||------------------------------------------------------------------------------------------------||
//|| Boolean to int8 pointer
//||------------------------------------------------------------------------------------------------||

func ToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

//||------------------------------------------------------------------------------------------------||
//|| Boolean to int8 pointer
//||------------------------------------------------------------------------------------------------||

func BoolToInt8Ptr(b bool) *int8 {
	var i int8 = 0
	if b {
		i = 1
	}
	return &i
}

//||------------------------------------------------------------------------------------------------||
//|| DerefInt8
//||------------------------------------------------------------------------------------------------||

func DerefInt8(i *int8) int8 {
	if i == nil {
		return 0
	}
	return *i
}

//||------------------------------------------------------------------------------------------------||
//|| IP to uint32
//||------------------------------------------------------------------------------------------------||

func IpToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}
