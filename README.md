# logging

Reusable Go logging package extracted from your `user_service`.

## Module

Current module path:

`github.com/learncodexx/logging`

If your GitHub repo uses a different path, update `go.mod`.

## Install

```bash
go get github.com/learncodexx/logging
```

## Basic usage

```go
import "github.com/learncodexx/logging"

basic := logging.NewBasicPrint()
basic.Info("START", "service up")
```

## API usage (Fiber)

```go
import "github.com/learncodexx/logging"

api := logging.NewAPIPrint("UserService")
logging.SetFiberErrorHook(middleware.AddErrorToRequest)

ctx, cancel := logging.WithDefaultTimeout(c, "UserService")
defer cancel()

if err := svc.SignIn(ctx, req); err != nil {
	return logging.Trace(err)
}

api.Info(ctx, "signin success")
```

## HTTP Middleware Usage

```go
import (
	"github.com/gofiber/fiber/v2"
	"github.com/learncodexx/logging"
)

func main() {
	app := fiber.New()

	// Enable grouped request logging.
	app.Use(logging.HTTPLogger())

	// Let APIPrint errors join the current HTTP log group.
	logging.SetFiberErrorHook(logging.AddErrorToRequest)
}
```

## Main functions

- `Trace(err error) error`
- `SetError(message string) error`
- `ErrorPattern(err error) string`
- `TranslateError(raw string) string`
- `NewBasicPrint() *BasicPrint`
- `NewAPIPrint(service string) *APIPrint`
- `SetServiceName(ctx, service)`
- `WithTimeout(...)`, `WithDefaultTimeout(...)`
- `SetFiberErrorHook(fn)`
