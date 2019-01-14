package dbutil

import (
	"sort"
	"testing"
)

func TestInt64ArrayType(t *testing.T) {
	var ary Int64Array
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	ary.Value()
	//
	v0, v1, v2 := int64(3), int64(2), int64(1)
	ary = append(ary, &v0)
	ary = append(ary, &v1)
	ary = append(ary, &v2)
	sort.Sort(ary)
}

func TestMap(t *testing.T) {
	var ary Map
	ary.Value()
	ary = Map{}
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("xxx")
	if err == nil {
		t.Error(err)
		return
	}

	ary.Value()
}

func TestStringArrayType(t *testing.T) {
	var ary StringArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}

	err = ary.Scan("xxx")
	if err == nil {
		t.Error(err)
		return
	}

	ary.Value()
	ary = StringArray{"xxx"}
	ary.Value()
}
