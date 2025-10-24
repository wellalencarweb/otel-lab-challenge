package customerrors

type NotFoundError struct {
	Err     error
	Message string
	Tags    map[string]interface{}
}

func (e *NotFoundError) Error() string {
	return e.Message
}
