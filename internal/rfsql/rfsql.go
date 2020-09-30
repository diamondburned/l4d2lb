package rfsql

import (
	"reflect"
	"sync"
)

var cache sync.Map

func Columns(v interface{}) []string {
	var t = reflect.TypeOf(v)
	if columns, ok := cache.Load(t); ok {
		return columns.([]string)
	}

	var n = t.NumField()
	var cols = make([]string, 0, n)

	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue // skip unexported fields
		}

		if db, ok := t.Field(i).Tag.Lookup("db"); ok {
			cols = append(cols, db)
		}
	}

	cache.Store(t, cols)
	return cols
}
