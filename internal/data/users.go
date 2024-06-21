package data

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Users struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" form:"name"`
	Nickname string `json:"nickname" form:"nickname"`
	Email    string `json:"email" form:"email"`
	City     string `json:"city" form:"city"`
	About    string `json:"about" form:"about"`
	Image    string `json:"image"`
}

var (
	mandatoryFieldRuLang = "Это поле является обязательным"
)

func (u Users) ValidateUsers() error {
	return validation.Errors{
		"Name": validation.Validate(u.Name,
			validation.Required.Error(mandatoryFieldRuLang),
			validation.Length(2, 20).Error("Имя должно содержать от 5 до 20 символов"),
		),

		"Nickname": validation.Validate(u.Nickname,
			validation.Required.Error(mandatoryFieldRuLang),
		),

		"Email": validation.Validate(u.Email,
			validation.Required.Error(mandatoryFieldRuLang),
			is.Email.Error("Должен быть действительный адрес электронной почты"),
		),

		"City": validation.Validate(u.City,
			validation.Required.Error(mandatoryFieldRuLang),
			validation.Match(regexp.MustCompile("^[a-zA-Z]+$")),
		),

		"About": validation.Validate(u.About,
			validation.Required.Error(mandatoryFieldRuLang),
			validation.Length(10, 200),
		),
	}.Filter()
}
