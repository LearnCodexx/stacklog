package stacklog

import "context"

// SetServiceName stores the service name in context for APIPrint tags.
func SetServiceName(c context.Context, name string) context.Context {
	return context.WithValue(c, KeyAPIPrint, name)
}
