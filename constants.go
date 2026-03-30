package stacklog

const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"

	TagAPI = "API"

	// Context key for overriding service/tag name when printing.
	KeyAPIPrint = "stacklog_service"
)
