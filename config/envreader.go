package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

var getEnv = func(key string) string {
	return os.Getenv(key)
}

// env2Map accepts any struct interface and uses the mapstructure tag to read values from environment variables
func env2Map(v interface{}) (*map[string]interface{}, error) {
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		return nil, errors.New("cannot decode from pointer")
	}

	m := map[string]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(TagName)
		val := getEnv(tag)

		switch field.Type.Kind() {
		case reflect.String:
			m[tag] = val
		case reflect.Int:
			d, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			m[tag] = d
		}
	}
	if err := mapstructure.Decode(&m, &v); err != nil {
		return nil, err
	}
	return &m, nil
}
