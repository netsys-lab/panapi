package taps

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func get(st interface{}, key string) (value reflect.Value, err error) {
	s := reflect.ValueOf(st).Elem()
	reg := regexp.MustCompile("[^a-z]+")
	stripKey := reg.ReplaceAllString(strings.ToLower(key), "")
	f := s.FieldByNameFunc(func(k string) bool {
		if reg.ReplaceAllString(strings.ToLower(k), "") == stripKey {
			return true
		} else {
			return false
		}
	})
	if !f.IsValid() {
		return reflect.ValueOf(nil), fmt.Errorf("Type %T has no Field %s (%s)", st, key, stripKey)
	}
	return f, nil
}

func set(st interface{}, key string, value interface{}) error {
	f, err := get(st, key)
	if err != nil {
		return err
	}
	p := reflect.ValueOf(value)
	if p.IsValid() && p.Type().AssignableTo(f.Type()) {
		f.Set(p)
	} else {
		return fmt.Errorf("Can not assign value of Type %T to Field %s of Type %T (expect %s)", value, key, st, f.Type())
	}
	return nil

}
