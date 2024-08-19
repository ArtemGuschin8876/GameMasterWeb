package data

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type TemplateData struct {
	Errors     []string          `json:"errors,omitempty"`
	FormErrors map[string]string `json:"form_errors,omitempty"`
	User       *User             `json:"user,omitempty"`
	Flash      string            `json:"flash,omitempty"`
	Users      []User            `json:"users,omitempty"`
}

type JsonData struct {
	FormErrors map[string]string `json:"form_errors"`
}

func (t TemplateData) ValidationOfFormFields(data TemplateData, err error) {
	if val, ok := err.(validation.Errors); ok {
		for field, valerr := range val {
			switch field {
			case "Name":
				data.FormErrors["name"] = valerr.Error()
			case "Nickname":
				data.FormErrors["nickname"] = valerr.Error()
			case "Email":
				data.FormErrors["email"] = valerr.Error()
			case "City":
				data.FormErrors["city"] = valerr.Error()
			case "About":
				data.FormErrors["about"] = valerr.Error()
			}
		}
	}
}
