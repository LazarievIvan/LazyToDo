package models

// DBError defines custom error for specific handling and chaining.
type DBError struct {
	message string
	code    int
	err     error
}

func NewDBError(msg string, code int, err error) *DBError {
	return &DBError{
		message: msg,
		code:    code,
		err:     err,
	}
}

func (e *DBError) Unwrap() error {
	return e.err
}

func (e *DBError) Error() string {
	return e.message
}

func (e *DBError) Code() int {
	return e.code
}
