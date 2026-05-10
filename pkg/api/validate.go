package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeAndValidate decodes a JSON request body into the given struct
// and validates it using `validate` struct tags. Returns a user-facing
// error message on failure, or nil on success.
func DecodeAndValidate(body io.ReadCloser, req interface{}) error {
	if err := json.NewDecoder(body).Decode(req); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %s", err.Error())
	}
	return nil
}
