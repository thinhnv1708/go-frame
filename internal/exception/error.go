package exception

type AppError struct {
	Code       int
	Message    string
	Err        error
	HttpStatus int
}

func (e *AppError) Error() string {
	return e.Message
}
