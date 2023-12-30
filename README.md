# RapidRoot

RapidRoot is a lightweight and flexible HTTP router library for Go. It provides a simple API to define routes, handle middleware, and manage HTTP requests and responses.

## Getting Started

To use RapidRoot in your Go project, follow these steps:

1. Install RapidRoot using `go get`:

   ```bash
   go get github.com/Folium1/RapidRoot
   ```

2. Import RapidRoot in your code:

   ```go
   import (
       "fmt"
       "net/http"
       "strings"
       "github.com/Folium1/RapidRoot"
   )
   ```

3. Create a new router instance:

   ```go
   r := RapidRoot.NewRouter()
   ```

4. Define routes using HTTP methods:

   ```go
   r.GET("/", func(req *RapidRoot.Request) {
       req.JSON(http.StatusOK, map[string]string{"message": "Hello, RapidRoot!"})
   })

   r.POST("/users", func(req *RapidRoot.Request) {
       // Handle POST request for creating users
   })
   ```

5. Run the HTTP server:

   ```go
   r.Run(":8080")
   ```

## Examples

### Basic Routing

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/Folium1/RapidRoot"
)

func main() {
    r := RapidRoot.NewRouter()

    // Define a simple GET route
    r.GET("/", func(req *RapidRoot.Request) {
        req.JSON(http.StatusOK, map[string]string{"message": "Hello, RapidRoot!"})
    })

    // Define a POST route
    r.POST("/users", func(req *RapidRoot.Request) {
        // Handle POST request for creating users
    })

    // Run the server on port 8080
    r.Run(":8080")
}
```

### Middleware

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/Folium1/RapidRoot"
)

func middleware1(next RapidRoot.HandlerFunc) RapidRoot.HandlerFunc {
    return func(req *RapidRoot.Request) {
        fmt.Println("Executing Middleware 1")
        next(req)
    }
}

func middleware2(next RapidRoot.HandlerFunc) RapidRoot.HandlerFunc {
    return func(req *RapidRoot.Request) {
        fmt.Println("Executing Middleware 2")
        next(req)
    }
}

func main() {
    r := RapidRoot.NewRouter()

    // Apply middleware to all routes
    r.MIDDLEWARE("/", middleware1, middleware2)

    // Define a GET route
    r.GET("/", func(req *RapidRoot.Request) {
        req.JSON(http.StatusOK, map[string]string{"message": "Hello, Middleware!"})
    })

    // Run the server on port 8080
    r.Run(":8080")
}
```
## Logging

RapidRoot provides a simple logging mechanism that outputs logs to the standard output. You can customize the log output by using the `SetOutput` function.

### SetOutput(w io.Writer)

Sets the output writer for the logger.

## Contributing

If you find any issues or have suggestions for improvement, feel free to open an issue or create a pull request on the [GitHub repository](https://github.com/Folium1/RapidRoot).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.