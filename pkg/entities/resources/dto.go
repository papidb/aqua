package resources

import validation "github.com/go-ozzo/ozzo-validation"

// UpdateResourceDTO
type UpdateResourceDTO struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (r UpdateResourceDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&r.Region, validation.Required, validation.Length(2, 50)),
	)
}
