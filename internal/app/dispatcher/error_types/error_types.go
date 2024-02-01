package error_types

import "errors"

var (
	ErrTaskAlreadyExist = errors.New("task already exists")
)
