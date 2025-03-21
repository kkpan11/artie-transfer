package maputil

import (
	"fmt"
	"strconv"

	"github.com/artie-labs/transfer/lib/typing"
)

func GetKeyFromMap(obj map[string]any, key string, defaultValue any) any {
	if len(obj) == 0 {
		return defaultValue
	}

	val, isOk := obj[key]
	if !isOk {
		return defaultValue
	}

	return val
}

func GetInt32FromMap(obj map[string]any, key string) (int32, error) {
	if len(obj) == 0 {
		return 0, fmt.Errorf("object is empty")
	}

	valInterface, isOk := obj[key]
	if !isOk {
		return 0, fmt.Errorf("key: %s does not exist in object", key)
	}

	val, err := strconv.ParseInt(fmt.Sprint(valInterface), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("key: %s is not type integer: %w", key, err)
	}

	return int32(val), nil
}

func GetTypeFromMap[T any](obj map[string]any, key string) (T, error) {
	value, isOk := obj[key]
	if !isOk {
		var zero T
		return zero, fmt.Errorf("key: %q does not exist in object", key)
	}

	return typing.AssertType[T](value)
}
