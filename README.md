# Logging package

One import for all logging helpers used by this service.

## Install / import
```go
import "learncodexx/point_of_sale/user_service/pkg/logging"
```

## Quick start
```go
log := logging.NewAPIPrint("user-service")
ctx := logging.SetServiceName(context.Background(), "user-service")

if err := doWork(); err != nil {
    return logging.Trace(err) // adds [file:line] stack hints
}

log.Info(ctx, "processed %d users", 5)
```

## API surface
- `Trace(err error) error` — wrap errors with caller file:line, skipping duplicates.
- `SetError(msg string) error` — create a fresh error with caller file:line.
- `APIPrint` — request-aware logger; groups errors with Fiber requests when `SetFiberErrorHook` is set.
- `BasicPrint` — lightweight stdout logger for init/background code.
- `CheckType(args ...any) string` — appends `[val]` hints to log lines.
- `ErrorPattern(err) string` / `TranslateError(raw string)` — convert DB/infra errors to user-friendly text.
- Context helpers: `WithTimeout`, `WithDefaultTimeout`, `BackgroundWithDefaultTimeout(Basic)`, `SetServiceName`.
- `SetFiberErrorHook(fn)` — connect Fiber middleware to group error logs per request.

## Patterns
- Call `logging.Trace(err)` inline at the return site (no defers) for accurate stack lines.
- Use `SetServiceName` or the `With*Timeout` helpers so APIPrint tags logs with the current service.
- For new errors, prefer `SetError` over `fmt.Errorf` to keep format consistent.
