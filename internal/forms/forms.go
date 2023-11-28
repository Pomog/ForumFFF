package forms

import (
	"net/http"
	"net/url"
	"strings"
)


type Form struct{
	url.Values
	Errors errors
}
func (f *Form)Valid()bool{
	return len(f.Errors)==0
}

func NewForm(data url.Values) *Form{
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form) Required(fields ...string){
	for _,field:=range fields{
		value:=f.Get(field)
		if strings.TrimSpace(value) == ""{
			f.Errors.Add(field,"This field is required!")
		}
	}
}

func (f *Form)Has(field string, r *http.Request)bool{
	dataFromHtml:=r.FormValue(field)
	if len(dataFromHtml)==0{
		return false
	}
	return true
}