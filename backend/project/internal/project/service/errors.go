package service

import "errors"

var (
	ErrNotFound        = errors.New("project not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUnauthenticated = errors.New("x-user-id is required")
)
