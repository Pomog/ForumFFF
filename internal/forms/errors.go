package forms

type errors map[string][]string

// Add adds and error message for a given form field
func (e errors) Add(field, messagae string) {
	e[field] = append(e[field], messagae)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 {
		return ""
	}
	return errorString[0]
}

// errorKey is a custom type for the context key
type ErrorKey string
