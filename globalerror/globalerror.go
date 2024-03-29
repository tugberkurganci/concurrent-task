package globalerror

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ErrorResponse struct {
	Status      int32                 `json:"status"`
	ErrorDetail []ErrorResponseDetail `json:"errorDetail"`
}

type ErrorResponseDetail struct {
	FieldName   string `json:"fieldName"`
	Description string `json:"description"`
}

var ValidationErrorDescriptionMap = map[string]string{
	"min":       "Your value should be greater than ",
	"required":  "Your value is mandatory",
	"acceptAge": "Your value should be greater than 18",
}

type CustomValidationError struct {
	HasError bool
	Field    string
	Tag      string
	Param    string
	Value    interface{}
}

func Validate(data interface{}) []CustomValidationError {
	var customValidationError []CustomValidationError

	if errors := validate.Struct(data); errors != nil {
		for _, fieldError := range errors.(validator.ValidationErrors) {
			var cve CustomValidationError
			cve.HasError = true
			cve.Field = fieldError.Field()
			cve.Tag = fieldError.Tag()
			cve.Param = fieldError.Param()
			cve.Value = fieldError.Value()
			customValidationError = append(customValidationError, cve)
		}
	}

	return customValidationError
}

func HandleValidationErrors(c *fiber.Ctx, errors []CustomValidationError) error {
	var errorResponse ErrorResponse
	var errorDetailList []ErrorResponseDetail

	for _, validationError := range errors {
		var errorDetail ErrorResponseDetail
		errorDetail.FieldName = validationError.Field
		errorDetail.Description = fmt.Sprintf("%s field has an error because %s%s", validationError.Field, ValidationErrorDescriptionMap[validationError.Tag], validationError.Param)
		errorDetailList = append(errorDetailList, errorDetail)
	}
	errorResponse.Status = http.StatusBadRequest
	errorResponse.ErrorDetail = errorDetailList

	return c.Status(http.StatusBadRequest).JSON(errorResponse)
}
