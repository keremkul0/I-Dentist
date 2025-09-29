package errors

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
