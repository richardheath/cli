package validators

import (
	"errors"
	"strconv"
)

// Required Check if rawValue is an empty string.
func Required(rawValue string) error {
	if rawValue == "" {
		return errors.New("Required.")
	}

	return nil
}

// Int Check if rawValue is a valid int.
func Int(rawValue string) error {
	if _, err := strconv.ParseInt(rawValue, 10, 64); err != nil {
		return errors.New("Must be a valid int.")
	}

	return nil
}

// Number Check if rawValue is a valid number.
func Number(rawValue string) error {
	if _, err := strconv.ParseFloat(rawValue, 32); err != nil {
		return errors.New("Must be a valid number.")
	}

	return nil
}

// Bool Check if rawValue is a valid boolean value.
func Bool(rawValue string) error {
	if _, err := strconv.ParseFloat(rawValue, 32); err != nil {
		return errors.New("Must be a valid boolean value.")
	}

	return nil
}
