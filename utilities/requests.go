package utilities

import (
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

var Validate *validator.Validate

func ValidateInterface(obj interface{}) error {
	err := Validate.Struct(obj)

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println(err)
		return nil
	}

	if _, ok := err.(*validator.ValidationErrors); ok {
		fmt.Println(err)
		return err
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

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
