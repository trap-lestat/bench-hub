package service

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrStopped            = errors.New("stopped")
	ErrInvalidScriptType  = errors.New("invalid script type")
	ErrUnsupportedEngine  = errors.New("unsupported engine")
)
