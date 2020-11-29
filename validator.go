package main

import (
	"reflect"
	"strconv"
	"strings"
)

// ValidateAllFields validates whether all fields of a struct are valid
func ValidateAllFields(item interface{}) (bool, error) {
	v := reflect.ValueOf(item)
	for i := 0; i < v.NumField(); i++ {
		validator := getValidatorFromValue(v.Field(i))
		if validator != nil {
			res, err := validator.validate(v.Type().Field(i), v.Field(i))
			if err != nil || res == false {
				return res, err
			}
		}
	}

	return true, nil
}

func getValidatorFromValue(fv reflect.Value) validator {
	switch fv.Interface().(type) {
	case int:
		return createIntValidator("validate")
	default:
		return nil
	}
}

// validator validates a struct field
type validator interface {
	validate(f reflect.StructField, fv reflect.Value) (bool, error)
}

// intValidator validates int fields of a strcut
type intValidator struct {
	tagName string
}

func createIntValidator(tag string) intValidator {
	return intValidator{
		tagName: tag,
	}
}

// validate validates whether an integer is valid
func (v intValidator) validate(f reflect.StructField, fv reflect.Value) (bool, error) {
	tag := f.Tag.Get(v.tagName)

	args := strings.Split(tag, ";")

	for _, arg := range args {
		keyValue := strings.Split(arg, ":")
		key := keyValue[0]
		valueString := strings.Join(keyValue[1:], ":")

		if key == "min" {
			minValue, err := strconv.Atoi(valueString)
			if err != nil {
				return false, err
			}
			val := fv.Interface().(int)
			if val < minValue {
				return false, nil
			}
		}

		if key == "max" {
			maxValue, err := strconv.Atoi(valueString)
			if err != nil {
				return false, err
			}
			val := fv.Interface().(int)
			if val > maxValue {
				return false, nil
			}
		}
	}

	return true, nil
}
