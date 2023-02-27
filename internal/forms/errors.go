package forms

type errors map[string][]string

// Add adds an error message for a given form field
func (err errors) Add(field, message string) {
	err[field] = append(err[field], message)
}

// Get returns the first error message
func (err errors) Get(field string) string {
	errString := err[field]
	if len(errString) == 0 {
		return ""
	}
	return errString[0]
}
