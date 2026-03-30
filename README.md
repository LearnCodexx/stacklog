# stacklog (extracted from point_of_sale user_service)

Small logging helper with two pieces:

- `Trace(err error) error` — wraps errors with file:line, keeps nesting tidy, and skips runtime/internal frames.
- `APIPrint` / `BasicPrint` — colored console loggers with optional Fiber request grouping via `AddErrorToRequestFromMiddleware` hook.

## Install

```sh
go get github.com/learncodexx/stacklog@v0.1.0
```

(If you’re testing locally alongside `user_service`, add a replace in that module’s `go.mod`:
`replace github.com/learncodexx/stacklog => ../logging`)

## Quick start

```go
import (
    "context"
    "github.com/learncodexx/stacklog"
)

var log = stacklog.NewAPIPrint("user-service")

func handler(ctx context.Context) error {
    if err := doWork(); err != nil {
        return stacklog.Trace(err)
    }
    log.Info(ctx, "processed %d items", 5)
    return nil
}
```

## Notes
- Prefer calling `Trace(err)` inline at the return site. This is the most precise and lowest‑overhead pattern (no defers needed).
- If a lower layer already added a `[ file:line ]`, `Trace` skips adding another to avoid noise.
- `APIPrint.Error` will attempt to group with an HTTP request if you set `stacklog.AddErrorToRequestFromMiddleware` in your Fiber middleware.
- `ErrorPattern` and `TranslateError` are optional helpers to turn DB/network errors into user-friendly messages.
