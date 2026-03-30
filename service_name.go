package logging

import (
	"context"
	"learncodexx/point_of_sale/user_service/constants"
)

// SetServiceName stores the service name in context for downstream log tagging.
func SetServiceName(c context.Context, name string) context.Context {
	return context.WithValue(c, constants.KeyAPIPrint, name)
}
