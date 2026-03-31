package stacklog

import (
	"strings"
	"sync"
)

// ErrorMapping represents a configurable error translation rule
type ErrorMapping struct {
	Pattern     string
	Message     string
	CaseSensitive bool
}

// ErrorPatternRegistry manages configurable error translations
type ErrorPatternRegistry struct {
	mappings []ErrorMapping
	mu       sync.RWMutex
}

// Default error pattern registry
var defaultRegistry = NewErrorPatternRegistry()

// NewErrorPatternRegistry creates a new error pattern registry with default mappings
func NewErrorPatternRegistry() *ErrorPatternRegistry {
	registry := &ErrorPatternRegistry{
		mappings: []ErrorMapping{},
	}
	
	// Add default patterns
	registry.AddDefaultMappings()
	return registry
}

// AddMapping adds a new error translation mapping
func (r *ErrorPatternRegistry) AddMapping(pattern, message string, caseSensitive bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.mappings = append(r.mappings, ErrorMapping{
		Pattern:       pattern,
		Message:       message,
		CaseSensitive: caseSensitive,
	})
}

// AddDefaultMappings adds the standard error translation patterns
func (r *ErrorPatternRegistry) AddDefaultMappings() {
	// Database errors
	r.AddMapping("unique constraint", "The record already exists in our system.", false)
	r.AddMapping("users_email_key", "This email is already registered. Please use another email.", false)
	r.AddMapping("users_phone_key", "This phone number is already in use.", false)
	r.AddMapping("violates foreign key constraint", "This data cannot be deleted because it is being used by other records.", false)
	r.AddMapping("no rows in result set", "The requested data could not be found.", false)
	
	// Network/Connection errors
	r.AddMapping("connection refused", "Failed to connect to the server. Please check your internet connection.", false)
	r.AddMapping("connection reset", "Failed to connect to the server. Please check your internet connection.", false)
	r.AddMapping("context deadline exceeded", "The request timed out. Please try again.", false)
	r.AddMapping("timeout", "The request timed out. Please try again.", false)
	
	// Authentication/Authorization errors
	r.AddMapping("permission denied", "You don't have permission to perform this action.", false)
	r.AddMapping("forbidden", "You don't have permission to perform this action.", false)
	r.AddMapping("unauthorized", "Access denied. Please log in again.", false)
	r.AddMapping("token", "Access denied. Please log in again.", false)
	
	// Validation errors
	r.AddMapping("invalid argument", "Invalid request data.", false)
	r.AddMapping("bad request", "Invalid request data.", false)
	r.AddMapping("not found", "The requested data could not be found.", false)
}

// Translate finds the first matching pattern and returns the translated message
func (r *ErrorPatternRegistry) Translate(errStr string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	searchStr := errStr
	if len(r.mappings) > 0 {
		// Most mappings are case-insensitive, prepare lowercase version
		lowerError := strings.ToLower(errStr)
		
		for _, mapping := range r.mappings {
			var targetStr string
			var pattern string
			
			if mapping.CaseSensitive {
				targetStr = searchStr
				pattern = mapping.Pattern
			} else {
				targetStr = lowerError
				pattern = strings.ToLower(mapping.Pattern)
			}
			
			if strings.Contains(targetStr, pattern) {
				return mapping.Message
			}
		}
	}
	
	// Return original error if no pattern matches
	return errStr
}

// SetDefaultRegistry allows replacing the default error pattern registry
func SetDefaultRegistry(registry *ErrorPatternRegistry) {
	defaultRegistry = registry
}

// GetDefaultRegistry returns the current default error pattern registry
func GetDefaultRegistry() *ErrorPatternRegistry {
	return defaultRegistry
}

// AddErrorMapping is a convenience function to add mappings to the default registry
func AddErrorMapping(pattern, message string, caseSensitive bool) {
	defaultRegistry.AddMapping(pattern, message, caseSensitive)
}