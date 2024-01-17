# RapidRoot

## Overview

RapidRoot is a Go package offering efficient HTTP routing capabilities for web applications. It supports various HTTP methods, dynamic routing, middleware, and more.

## Features

- **HTTP Methods:** Supports GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE.
- **Middleware:** Route-specific and group middleware functionality.
- **Dynamic Routing:** Handles dynamic routes with path parameters.
- **Response Utilities:** Includes built-in methods for common HTTP responses (JSON, XML, HTML, etc.).
- **Request and Response Wrappers:** Enhances functionality and flexibility.
- **Cookie Management:** Secure and customizable handling of cookies.
- **Efficient Request Pooling:** Reduces garbage collection overhead.

## Installation

```bash
go get github.com/Folium1/RapidRoot
```

## Basic Usage

### Importing the Package

```go
import rr "github.com/your-username/rapidroot"
```

### Creating a Router

```go
router := rr.NewRouter()
```

### Defining Routes

```go
router.GET("/path", handlerFunction)
router.POST("/path", handlerFunction)
// Repeat for other HTTP methods
```

### Starting the Server

For HTTP:

```go
router.Run(":8080")
```

For HTTPS:

```go
router.RunWithTLS(":443", "certFile", "keyFile")
```

### Handler Function

```go
func handlerFunction(req *rr.Request) {
    // Request handling logic here
}
```

## Middleware
```go
func loggingMiddleware(next rapidroot.HandlerFunc) rapidroot.HandlerFunc {
    return func(req *rapidroot.Request) {
        log.Printf("Request received: %s %s", req.Req.Method, req.Req.URL.Path)
        next(req) // Call the next handler
    }
}

```

### Applying Middleware

```go
router.Middleware("GET", "/path", yourMiddlewareFunction)
```

### Group Middleware

```go
router.GroupMiddleware("GET", "/api", middlewareFunction1, middlewareFunction2)
```

## Advanced Features

- Custom request and response manipulation.
- Secure and flexible cookie handling.
- Dynamic routing with easy parameter extraction.

## Contributing

Contributions are welcome. Please adhere to Go's standard coding style and submit pull requests for any contributions.

## License

RapidRoot is released under the [MIT License](https://opensource.org/licenses/MIT).

---

*Note: For detailed API documentation and advanced usage, refer to the source code comments*