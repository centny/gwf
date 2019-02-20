package dbutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Centny/gwf/util"
)

//Time is database value to parse data from database and parset time.Time to timestamp on json mashal
type Time time.Time

//Timestamp return timestamp
func (t Time) Timestamp() int64 {
	return time.Time(t).Local().UnixNano() / 1e6
}

//MarshalJSON marshal time to string
func (t *Time) MarshalJSON() ([]byte, error) {
	raw := t.Timestamp()
	if raw < 0 {
		return []byte("null"), nil
	}
	stamp := fmt.Sprintf("%v", raw)
	return []byte(stamp), nil
}

//UnmarshalJSON unmarshal string to time
func (t *Time) UnmarshalJSON(bys []byte) (err error) {
	val := strings.TrimSpace(string(bys))
	if val == "null" {
		return
	}
	timestamp, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		*t = Time(time.Unix(0, timestamp*1e6))
	}
	return
}

//Scan is sql.Sanner
func (t *Time) Scan(src interface{}) (err error) {
	if src != nil {
		if timeSrc, ok := src.(time.Time); ok {
			*t = Time(timeSrc)
		}
	}
	return
}

//Int64Array is database value to parse data to []int64 value
type Int64Array []*int64

//Scan is sql.Sanner
func (i *Int64Array) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (i *Int64Array) Value() string {
	if i == nil || *i == nil {
		return "[]"
	}
	bys, _ := json.Marshal(*i)
	return string(bys)
}

//DbJoin will parset to database array
func (i Int64Array) DbJoin() string {
	vals := []string{}
	for _, v := range i {
		if v == nil {
			vals = append(vals, "nil")
		} else {
			vals = append(vals, fmt.Sprintf("%v", *v))
		}
	}
	return strings.Join(vals, ",")
}

func (i Int64Array) Len() int {
	return len(i)
}
func (i Int64Array) Less(a, b int) bool {
	return *i[a] < *i[b]
}
func (i Int64Array) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i Int64Array) HavingOne(vals ...int64) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//IntArray is database value to parse data to []int value
type IntArray []*int

//Scan is sql.Sanner
func (i *IntArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), i)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (i *IntArray) Value() string {
	if i == nil || *i == nil {
		return "[]"
	}
	bys, _ := json.Marshal(*i)
	return string(bys)
}

//DbJoin will parset to database array
func (i IntArray) DbJoin() string {
	vals := []string{}
	for _, v := range i {
		vals = append(vals, fmt.Sprintf("%v", *v))
	}
	return strings.Join(vals, ",")
}

func (i IntArray) Len() int {
	return len(i)
}
func (i IntArray) Less(a, b int) bool {
	return *i[a] < *i[b]
}
func (i IntArray) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i IntArray) HavingOne(vals ...int) bool {
	for _, v0 := range i {
		for _, v1 := range vals {
			if *v0 == v1 {
				return true
			}
		}
	}
	return false
}

//Map is database value to parse json data to map value
type Map util.Map

//Scan is sql.Sanner
func (m *Map) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), m)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (m *Map) Value() string {
	if m == nil || *m == nil {
		return "{}"
	}
	bys, _ := json.Marshal(*m)
	return string(bys)
}

//StringArray is database value to parse data to []string value
type StringArray []string

//Scan is sql.Sanner
func (s *StringArray) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), s)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (s *StringArray) Value() string {
	if s == nil || *s == nil {
		return "[]"
	}
	bys, _ := json.Marshal(*s)
	return string(bys)
}
