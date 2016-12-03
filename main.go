// under $GOPATH/{src}/{host}/{authorname}/{projectname}
// so it can work with `go get ...`
package main

// dependency
import "net/http"

func main() {
    // HandlerFunc register "/" with hello for ServeMux(http router)
    // "/" matches all path, not only the "/"
    // first-class function
    // [!] So can see this function change and save some internal states we cannot see here
    http.HandleFunc("/", hello)     // (pattern, handler that been registered)
    http.ListenAndServe(":8080", nil)
}

// with a type system of strong, static, inferred
// http.HandlerFunc
// server will spawn a new goroutine executing this for each request
func hello(w http.ResponseWriter, r *http.Request) {
    // w.Write take general []byte, so need to cast
    // if lost ')', compiler will raise: ./main.go:19: syntax error: unexpected semicolon or newline, expecting comma or )
    // it is ok
    w.Write([]byte("hello!"))
}
