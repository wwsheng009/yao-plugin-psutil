package main

import (
	"fmt"
	"reflect"
)

const (
	b  = uint64(1)
	kb = 1024 * b
	mb = 1024 * kb
	gb = 1024 * mb
)

func formatBytes(bytes uint64) string {
	switch {
	case bytes < kb:
		return fmt.Sprintf("%dB", bytes)
	case bytes < mb:
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(kb))
	case bytes < gb:
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(mb))
	default:
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(gb))
	}
}

// copy data to data2,and convert the fields
// data and data2 should have the similar struct
func convertToFormattedString(data interface{}, data2 interface{}) {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	v2 := reflect.ValueOf(data2)
	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		field2 := v2.Field(i)
		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			elem := field.Elem()
			// elem2 := field2.Elem()

			if elem.Kind() == reflect.Struct {
				field2.Set(reflect.New(field2.Type().Elem()))
				elem2 := field2.Elem()
				convertToFormattedString(elem.Addr().Interface(), elem2.Addr().Interface())
			}
		case reflect.Slice:
			sliceType := reflect.SliceOf(field2.Type().Elem())
			newSlice := reflect.MakeSlice(sliceType, field.Len(), field.Len())
			field2.Set(newSlice)

			for j := 0; j < field.Len(); j++ {
				convertToFormattedString(field.Index(j).Addr().Interface(), field2.Index(j).Addr().Interface())
			}
		case reflect.Uint64, reflect.Int:
			str := formatBytes(field.Uint())
			field2.SetString(str)
		case reflect.Struct:
			convertToFormattedString(field.Addr().Interface(), field2.Addr().Interface())
		default:
			field2.Set(field)
		}
	}
}
