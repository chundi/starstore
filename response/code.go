package response


const (
	OK int = 0

	Error int = 40000
	NotFound = iota

	DBConnectionError = iota
)