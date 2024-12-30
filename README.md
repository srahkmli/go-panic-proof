# panicrecovery

`panicrecovery` is a Go package that provides middleware for recovering from panics in both HTTP and gRPC servers. This middleware helps ensure that  server doesn't crash when an unexpected panic occurs, logging the error and providing an appropriate response.

## Features

- Recover from panics in HTTP handlers.
- Recover from panics in gRPC service methods.
- Logs the panic details along with the stack trace for debugging.
- Responds with `500 Internal Server Error` for HTTP, ensuring clients are informed about the issue.
```bash
2024/12/18 15:32:45 Recovered from panic in gRPC method /service/SomeMethod: something went wrong!
Stack Trace: 
goroutine 10 [running]:
runtime/debug.Stack(0x0, 0x0, 0x0)
	/usr/local/go/src/runtime/debug/stack.go:24 +0x88
github.com/srahkmli/panicrecovery.RecoverInterceptor.func1(0x1f7cae0, 0x0, 0x1f7cae0, 0x0, 0x7fddc35d19b0)
	/project/panicrecovery/grpc.go:23 +0xd3
google.golang.org/grpc.(*Server).processUnaryRPC(0xc0004db5e0, 0x7fddc35d1928, 0x1f9b960, 0xc0004d7030, 0xc00023bc00, 0x1f7cb80, 0x0, 0x0, 0x0)
	/project/go/pkg/mod/google.golang.org/grpc@v1.42.0/server.go:1281 +0x45b
...

```
  
## Installation

To install the `panicrecovery` package, run the following command in  Go project:

```bash
go get github.com/srahkmli/panicrecovery
``` 
## Overview

### `HTTPRecover`

- **Description:** A simple middleware that recovers from panics, logs the error, and sends a default `500 Internal Server Error` response.
- **Key Features:**
  - Logs the panic error and stack trace.
  - Responds with a fixed error message (`"Internal Server Error"`) and HTTP status code (`500`).
- **Use Case:** Suitable for applications requiring a straightforward panic recovery mechanism without customization.

#### Example Usage:
```go
http.Handle("/example", HTTPRecover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    panic("example panic")
})))
```

---

### `HTTPRecoverWithHandler`

- **Description:** A flexible middleware that allows custom handling of errors via a user-provided `HTTPErrorHandler` function.
- **Key Features:**
  - Accepts an `HTTPErrorHandler` function for customizable error responses.
  - Provides a fallback behavior similar to `HTTPRecover` if no handler is specified.
- **Use Case:** Ideal for applications requiring detailed and tailored error responses or additional error-handling logic, such as sending alerts or metrics.

#### Example Usage:
```go
customHandler := func(w http.ResponseWriter, r *http.Request, err interface{}) {
    log.Printf("Custom handler: panic %v", err)
    http.Error(w, "Something went wrong", http.StatusInternalServerError)
}

http.Handle("/example", HTTPRecoverWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    panic("example panic")
}), customHandler))
```

- Use **`HTTPRecover`** for simple applications that do not require custom error handling.

- Use **`HTTPRecoverWithHandler`** when you need to provide detailed error responses or additional error-handling logic. The additional flexibility ensures the middleware can adapt to complex application requirements.
---

## Differences

| Aspect                | `HTTPRecover`                                           | `HTTPRecoverWithHandler`                                  |
|-----------------------|---------------------------------------------------------|----------------------------------------------------------|
| **Error Handling**     | Fixed behavior: log + `500 Internal Server Error`.      | Customizable via `HTTPErrorHandler`.                     |
| **Flexibility**        | Minimal flexibility.                                    | Highly flexible due to user-provided handler.            |
| **Custom Error Output**| Not supported.                                          | Fully supported via the handler.                        |
| **Default Fallback**   | Logs the error and responds with `500`.                 | Logs the error and responds with `500` if no handler.    |

---
 


