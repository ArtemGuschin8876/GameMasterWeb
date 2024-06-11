package data

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Users struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	City     string `json:"city"`
	About    string `json:"about"`
	Image    string `json:"image"`
}

func (u Users) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.Required,
			validation.Length(2, 20),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$"))),

		validation.Field(&u.Nickname,
			validation.Required,
			validation.Length(2, 20),
		),

		validation.Field(&u.Email,
			validation.Required,
			validation.Length(1, 100),
			is.Email),

		validation.Field(&u.City,
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")),
			validation.Length(1, 20),
			validation.Required),

		validation.Field(&u.About,
			validation.Required,
			validation.Length(30, 200),
		),
	)
}
