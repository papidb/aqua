package customers

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateCustomerDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validation rules for the struct
func (r CreateCustomerDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}

type ErrExistingEmailOrName struct{}

func (e ErrExistingEmailOrName) Error() string {
	return "email or name already in use"
}

type ErrExistingCustomerResource struct{}

func (e ErrExistingCustomerResource) Error() string {
	return "This relationship already exists"
}

type AddResourceToCustomerDTO struct {
	ResourceID string `json:"resource_id"`
}

func (r AddResourceToCustomerDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ResourceID, validation.Required, is.UUID),
	)
}
