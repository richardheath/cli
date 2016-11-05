package converters

import "strconv"

// Int convert rawValue to int.
func Int(rawValue string) (interface{}, error) {
	return strconv.Atoi(rawValue)
}

// Bool convert rawValue to bool.
func Bool(rawValue string) (interface{}, error) {
	return strconv.ParseBool(rawValue)
}

// Float convert rawValue to float.
func Float(rawValue string) (interface{}, error) {
	return strconv.ParseFloat(rawValue, 32)
}
