package customerrors

type ValidationError struct {
	Err     error
	Message string
	Reasons []string
	Tags    map[string]interface{}
}

func (e *ValidationError) Error() string {
	return e.Message
}
