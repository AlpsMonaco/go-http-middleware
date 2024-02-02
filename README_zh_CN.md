# go-http-middleware
一个用于go的`net/http`标准库的中间件库  
为`net/http`库提供便捷的中间库集成方案  

# Features
* 不需要依赖其他第三方库
* 代码简单易懂
* 和 `net/http` 库 完全兼容


# 快速开始
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
编写你自己的中间件，然后添加以下两行代码到你注册http handlers的地方  
```go
http := middleware.DefaultHTTPBuilder()
http.With(LogMiddleware1, LogMiddleware2)
```
没有其他地方需要修改

## 测试
```bash
curl http://127.0.0.1:33333/hello
```
控制台会输出:  
```
2024/02/02 14:00:54 log middleware 1 in
2024/02/02 14:00:54 log middleware 2 in
2024/02/02 14:00:54 Hello
2024/02/02 14:00:54 log middleware 2 out
2024/02/02 14:00:54 log middleware 1 out
```

# 实用中间件
## 日志
记录地址和路径
```go
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.RemoteAddr, r.RequestURI)
		next(w, r)
	}
}
```

## 计时器
```go
func TimerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		next(w, r)
		fmt.Printf("time elapsed:%d\nms",time.Now().Sub(st).Milliseconds())
	}
}
```

## recover中间件
recover if panic
```go
func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.WriteHeader(500)
				fmt.Println(err)
			}
		}()
		next(w, r)
	}
}
```

# 进阶
如果你使用了自定义的 mux , 请使用`middleware.NewBuilder`

```go
...
b := middleware.NewBuilder(&http.ServeMux{})
b.With(LogMiddleware1, LogMiddleware2)
b.HandleFunc("/hello", HttpHandleHello)
b.HandleFunc("/world", HttpHandleWorld)
b.ListenAndServe(":33333", nil)
```