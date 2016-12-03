// under $GOPATH/{src}/{host}/{authorname}/{projectname}
// so it can work with `go get ...`
package main

// dependency
import (
	"encoding/json"
	"net/http"
	"strings"
)

// define a new type, which is a struct
type weatherData struct {
	Name string `json:"name"` // name type tag
	// tag helps to map json to this struct
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

// return (result, err) is an idiom of Go
func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=6c68ab3ee677cf2fe425ff3ac9ab9b6f&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	// `defer` will execute the function call before out of the scope of query function
	// mainly for resources releasing
	defer resp.Body.Close()

	// use `var` or directly `:=`
	var d weatherData

	// interface
	// json.NewDecoder takes io.Reader interface, not a concrete HTTP resp body
	// http.Response.Body satisfy io.Reader
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

func main() {
	// HandlerFunc register "/" with hello for ServeMux(http router)
	// "/" matches all path, not only the "/"
	// first-class function
	// [!] So can see this function change and save some internal states we cannot see here
	http.HandleFunc("/hello", hello) // (pattern, handler that been registered)

	// at first a mistype `HandleFunc` as `HandlerFunc` and the compiler
	// report it poorly: ./main.go:62: too many arguments to conversion to http.HandlerFunc: http.HandlerFunc("/weather/", func literal)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

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
