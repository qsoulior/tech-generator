package test_db

import (
	"reflect"
)

type field struct {
	key   string
	value any
}

func toFields[T any](item T) []field {
	v := reflect.ValueOf(item)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil
	}

	fields := make([]field, v.NumField())
	for i := range fields {
		fields[i].key = t.Field(i).Tag.Get("db")
		fields[i].value = v.Field(i).Interface()
	}

	return fields
}
