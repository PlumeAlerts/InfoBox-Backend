package utilities

import (
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

type Error struct {
	Message string
}

var Validate *validator.Validate

func ValidateInterface(obj interface{}) error {
	err := Validate.Struct(obj)

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return nil
	}

	if _, ok := err.(*validator.ValidationErrors); ok {
		return err
	}
	return nil
}

func GetAnnotationID(id string) (int, error) {
	err := Validate.Var(id, "required,numeric")
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return i, nil
}
