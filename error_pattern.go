package logging

import "strings"

// ErrorPattern turns a wrapped error into a short, user-facing message.
func ErrorPattern(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()
	parts := strings.Split(errStr, "->")
	if len(parts) < 2 {
		return TranslateError(errStr)
	}

	rawError := strings.TrimSpace(parts[len(parts)-1])
	rawError = strings.ReplaceAll(rawError, "ERROR:", "")
	rawError = strings.TrimSpace(rawError)

	return TranslateError(rawError)
}

// TranslateError maps raw DB/network strings into friendly messages.
func TranslateError(rawErr string) string {
	lowErr := strings.ToLower(rawErr)

	switch {
	case strings.Contains(lowErr, "unique constraint") || strings.Contains(lowErr, "duplicate key"):
		if strings.Contains(lowErr, "users_email_key") {
			return "This email is already registered. Please use another email."
		}
		if strings.Contains(lowErr, "users_phone_key") {
			return "This phone number is already in use."
		}
		return "The record already exists in our system."

	case strings.Contains(lowErr, "violates foreign key constraint"):
		return "This data cannot be deleted because it is being used by other records."

	case strings.Contains(lowErr, "permission denied") || strings.Contains(lowErr, "forbidden"):
		return "You don't have permission to perform this action."

	case strings.Contains(lowErr, "not found") || strings.Contains(lowErr, "no rows in result set"):
		return "The requested data could not be found."

	case strings.Contains(lowErr, "connection refused") || strings.Contains(lowErr, "connection reset"):
		return "Failed to connect to the server. Please check your internet connection."

	case strings.Contains(lowErr, "context deadline exceeded") || strings.Contains(lowErr, "timeout"):
		return "The request timed out. Please try again."

	case strings.Contains(lowErr, "unauthorized") || strings.Contains(lowErr, "token"):
		return "Access denied. Please log in again."

	case strings.Contains(lowErr, "invalid argument") || strings.Contains(lowErr, "bad request"):
		return "Invalid request data."

	case strings.Contains(lowErr, "of json input"):
		return "Invalid request body"

	case strings.Contains(lowErr, "unexpected eof") || strings.Contains(lowErr, "eof"):
		return "Request body is incomplete or malformed."

	default:
		idx := strings.LastIndex(lowErr, "]")
		if idx == -1 {
			return lowErr
		}

		lowErr = strings.TrimSpace(lowErr[idx+1:])

		if len(lowErr) > 25 {
			return "An internal system error occurred. Please try again later."
		}
		return lowErr
	}
}
