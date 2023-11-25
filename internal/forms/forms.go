package forms

import (
	"net/http"
	"net/url"
)


type Form struct{
	url.Values
	Errors errors
}

func NewForm(data url.Values) *Form{
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form)Has(field string, r *http.Request)bool{
	dataFromHtml:=r.FormValue(field)
	if len(dataFromHtml)==0{
		return false
	}
	return true
}