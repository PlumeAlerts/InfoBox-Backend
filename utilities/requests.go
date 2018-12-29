package utilities

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
)

var Validate *validator.Validate

func ValidateInterface(obj interface{}) validator.ValidationErrors {
	err := Validate.Struct(&obj)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}
	return nil
}

func GetIBID(id string) (uint, error) {
	err := Validate.Var(id, "required,numeric")
	if err != nil {
		return 0, err
	}

	t, _ := strconv.ParseUint(id, 10, 32)
	return uint(t), nil
}

func InterfaceToJson(obj interface{}) ([]byte, bool) {
	b, err := json.Marshal(&obj)
	if err != nil {
		return nil, true
	}
	return b, false
}
