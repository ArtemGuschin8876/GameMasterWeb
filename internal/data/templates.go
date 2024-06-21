package data

import validation "github.com/go-ozzo/ozzo-validation"

type TemplateData struct {
	Errors     []string
	FormErrors map[string]string
	User       Users
	Flash      string
	Users      []Users
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
