package response

const (
	OK int = iota
	Error
	NotFound
	Unauthorized

	MethodNotAllowed

	DBConnectionError = iota
)
