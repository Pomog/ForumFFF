package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Form is a custom type that embeds url.Values and includes an Errors field.
type Form struct {
	url.Values
	Errors errors
}

// Valid checks if there are any errors in the form data
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// NewForm creates a new form instance with provided data and an empty error map
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

// Required checks if specified fields are present and not empty in the form data
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required!")
		}
	}
}

// First_LastName_Min_Max_Len checks if the length of a field in the form data is within the specified range
func (f *Form) First_LastName_Min_Max_Len(field string, minLen, maxLen int, r *http.Request) bool {
	inputData := r.FormValue(field)
	if len(inputData) < minLen {
		f.Errors.Add(field, fmt.Sprintf("Too Short. %s field should be NLT %d characters long.", field, minLen))
		return false
	} else if len(inputData) > maxLen {
		f.Errors.Add(field, fmt.Sprintf("Too Long. %s field shoiuld be NMT %d characters long.", field, maxLen))
		return false
	}
	return true
}

// EmailFormat checks if the specified field in the form data matches a valid email format
func (f *Form) EmailFormat(field string, r *http.Request) bool {
	inputData := strings.ToLower(r.FormValue(field))
	if len(strings.Split(inputData, "@")) != 2 {
		f.Errors.Add(field, "Wrong format of Email - wrong number of @ sign")
		return false
	}
	sliceToCheck := strings.Split(inputData, "@")
	if len(sliceToCheck) == 2 {
		firstP := sliceToCheck[0]
		secondP := sliceToCheck[1]
		if len(firstP) == 0 {
			f.Errors.Add(field, "Wrong format of Email - should be local part before @, like john@, or nick@")
			return false
		}
		sliceToCheck2 := strings.Split(secondP, ".")
		if len(sliceToCheck2) == 2 {
			domName := sliceToCheck2[0]
			TLDName := sliceToCheck2[1]
			if len(domName) == 0 {
				f.Errors.Add(field, "Wrong format of Email - should be domain name after @, like @gmail, or @yahoo")
				return false
			} else if len(TLDName) == 0 {
				f.Errors.Add(field, "Wrong format of Email - should be TLD after @, like .com or .net")
				return false
			}
		}
	}
	return true
}

// PassFormat checks if the specified field in the form data matches certain password format criteria
func (f *Form) PassFormat(field string, minL, maxL int, r *http.Request) bool {
	inputData := r.FormValue(field)
	if len(inputData) < minL {
		f.Errors.Add(field, fmt.Sprintf("Too Short. %s field should be NLT %d characters long.", field, minL))
		return false
	} else if len(inputData) > maxL {
		f.Errors.Add(field, fmt.Sprintf("Too Long. %s field shoiuld be NMT %d characters long.", field, maxL))
		return false
	}
	return true
}

// Has checks if the specified field exists in the form data
func (f *Form) Has(field string, r *http.Request) bool {
	return len(r.FormValue(field)) != 0
}

func CheckSingleWordLen(input string, SingleWordMaxLen int) bool {
	splitted := strings.Split(input, " ")
	for _, elem := range splitted {
		if len(elem) > SingleWordMaxLen {
			return false
		}
	}
	return true

}
