package base_domain

type Error struct {
	Msg string
}

func NewError(msg string) *Error {
	return &Error{Msg: msg}
}

func (e *Error) Error() string { return e.Msg }
