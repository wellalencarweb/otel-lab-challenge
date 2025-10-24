package customerrors

type UnknownError struct {
	Err     error
	Message string
	Tags    map[string]interface{}
}

func (e *UnknownError) Error() string {
	return e.Message
}
