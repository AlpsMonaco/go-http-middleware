[中文文档](https://github.com/AlpsMonaco/go-http-middleware/blob/main/README_zh_CN.md)

# go-http-middleware
A middleware library for go standard net/http library.  

# Features
* no more third party library requires.
* core code are easy to learn and understand.
* fully compatible with the `net/http` library


# Quick Start
```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AlpsMonaco/go-http-middleware/middleware"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func HttpHandleHello(w http.ResponseWriter, r *http.Request) {
	logger.Println("Hello")
	w.Write([]byte("Hello"))
}

func HttpHandleWorld(w http.ResponseWriter, r *http.Request) {
	logger.Println("World")
	w.Write([]byte("World"))
}

func LogMiddleware1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("log middleware 1 in")
		next(w, r)
		logger.Println("log middleware 1 out")
	}
}

func LogMiddleware2(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("log middleware 2 in")
		next(w, r)
		logger.Println("log middleware 2 out")
	}
}

func main() {
	http := middleware.DefaultHTTPBuilder()
	http.With(LogMiddleware1, LogMiddleware2)
	http.HandleFunc("/hello", HttpHandleHello)
	http.HandleFunc("/world", HttpHandleWorld)
	http.ListenAndServe(":33333", nil)
}

```
Write your own middleware and add the following two lines of code before registering the HTTP handler in your code.  
```go
http := middleware.DefaultHTTPBuilder()
http.With(LogMiddleware1, LogMiddleware2)
```
No further modifications are needed.  

## Test
```bash
curl http://127.0.0.1:33333/hello
```
the console will output 
```
2024/02/02 14:00:54 log middleware 1 in
2024/02/02 14:00:54 log middleware 2 in
2024/02/02 14:00:54 Hello
2024/02/02 14:00:54 log middleware 2 out
2024/02/02 14:00:54 log middleware 1 out
```

# Advance
Use `middleware.NewBuilder` if you use your own mux instead of default mux of net/http.

```go
...
b := middleware.NewBuilder(&http.ServeMux{})
b.With(LogMiddleware1, LogMiddleware2)
b.HandleFunc("/hello", HttpHandleHello)
b.HandleFunc("/world", HttpHandleWorld)
b.ListenAndServe(":33333", nil)
```