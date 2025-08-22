package e

// Error codes for the application
const (
	SUCCESS        = 200
	ERROR          = 500
	InvalidParams  = 400
	Unauthorized   = 401
	NotFound       = 404
	Forbidden      = 403
)

// Custom error codes for user module (example)
const (
	ErrorUserNotFound      = 10001
	ErrorUserPasswordWrong = 10002
	ErrorUserAlreadyExists = 10003
)

// Add more modules' error codes here...
