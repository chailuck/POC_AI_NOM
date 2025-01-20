// internal/validation/validator.go
package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/your-username/tmf632-service/internal/models"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (cv *CustomValidator) ValidateIndividual(individual *models.Individual) error {
	// Add custom validation rules for Individual
	if individual.GivenName == "" {
		return fmt.Errorf("given name is required")
	}

	// Validate contact medium
	for _, cm := range individual.ContactMedium {
		if cm.Type == "PhoneContactMedium" && cm.PhoneNumber == "" {
			return fmt.Errorf("phone number is required for PhoneContactMedium")
		}
	}

	return nil
}
