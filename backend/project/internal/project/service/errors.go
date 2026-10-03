package service

import "errors"

var (
	ErrNotFound        = errors.New("project not found")
	ErrInvalidArgument = errors.New("invalid argument")
)
