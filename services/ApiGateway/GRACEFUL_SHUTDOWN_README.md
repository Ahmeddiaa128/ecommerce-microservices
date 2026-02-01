# Graceful Shutdown Tutorial (API Gateway)

This guide explains how the API Gateway implements graceful shutdown and how to validate it.

## Why graceful shutdown matters

Graceful shutdown:

- Stops accepting new connections.
- Lets in-flight requests finish within a timeout.
- Cancels request contexts so handlers can exit early.
- Prevents abrupt termination and lost responses.

## What the implementation does

1. Starts the HTTP server in a goroutine.
2. Listens for SIGINT/SIGTERM (or server errors).
3. Disables keep-alives to prevent new reuse of existing connections.
4. Cancels the server base context so handlers can observe shutdown.
5. Calls `Shutdown` with a timeout to let requests drain.
6. Waits for the server goroutine to finish.

## Key pieces to look at

- **Signal handling:** Uses a context created by `signal.NotifyContext`.
- **Server lifecycle:** Captures server errors through a channel.
- **Context cancellation:** Uses `BaseContext` with a cancellable root context.
- **Drain timeout:** Uses `context.WithTimeout` for `Shutdown`.

## How to test gracefully

1. Start the API Gateway.
2. Send a long-running request (e.g., a route that sleeps or streams).
3. Send SIGINT (Ctrl+C) or SIGTERM to the process.
4. Observe:
   - No new requests are accepted.
   - The in-flight request completes (within the timeout).
   - The process exits cleanly.

## Tips for production

- Set a timeout that matches your longest expected request.

### Ensure handlers respect context cancellation

- Always read `ctx := c.Request.Context()` (Gin) and pass it downstream.
- Before long work or loops, check:
  - `select { case <-ctx.Done(): return default: }`
- Ensure outbound calls (DB, gRPC, HTTP) use the same `ctx`.
- For streaming, stop writing when `ctx.Done()` is closed.

### Ensure background goroutines stop on shutdown

- Create goroutines with a parent `context.Context`.
- Use `select` on `<-ctx.Done()` inside goroutines and return when canceled.
- Keep a `sync.WaitGroup` and wait for all goroutines before exiting `main` if needed.
- Avoid orphaned goroutines started from handlers without a bounded lifetime.

### Keep logs clear and structured around shutdown steps

- Log the start of shutdown, the timeout value, and the reason (signal/server error).
- Log before and after `Shutdown`, and when all goroutines exit.
- Use consistent log keys (e.g., `component`, `event`, `timeout`).

## Related file

- services/ApiGateway/cmd/main.go
