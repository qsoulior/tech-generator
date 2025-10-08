package error_domain

type BaseError struct {
	Msg string
}

func NewBaseError(msg string) *BaseError {
	return &BaseError{Msg: msg}
}

func (e *BaseError) Error() string { return e.Msg }
