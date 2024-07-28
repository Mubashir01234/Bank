package errors

import "fmt"

type Error struct {
	ErrCode    string `json:"error_code,omitempty" example:"ERR000"`
	Err        error  `json:"-"`
	Message    string `json:"message" example:"Something went wrong"`
	StatusCode int    `json:"code" example:"500"`
}

func (e Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Err.Error()
}

func New(msg string, args ...interface{}) *Error {
	return &Error{StatusCode: 500, Message: fmt.Sprintf(msg, args...)}
}
