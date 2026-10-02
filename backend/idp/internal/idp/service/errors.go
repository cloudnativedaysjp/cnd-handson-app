package service

import "errors"

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
