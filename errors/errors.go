package errors

import "errors"

var (
	HostIsDown      = errors.New("hostIsDown")
	BadRequest      = errors.New("badRequest: ")
	InternalError   = errors.New("internal server error. please try again ")
	ConnectionIssue = errors.New("service is currently unable to respond. please try again")
)
