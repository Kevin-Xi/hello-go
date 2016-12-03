// under $GOPATH/{src}/{host}/{authorname}/{projectname}
// so it can work with `go get ...`
package main

// dependency
import (
    "encoding/json"
    "net/http"
    "strings"
    "log"
    "time"
)

// business logic
type weatherProvider interface {
    temperature(city string) (float64, error)
}

type openWeatherMap struct{
    apiKey string
}

// return (result, err) is an idiom of Go
// [!] this means it is duck typing?
// the `temperature` "menthod" of `openWeatherMap`
func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if err != nil {
		return 0, err
	}

	// `defer` will execute the function call before out of the scope of query function
	// mainly for resources releasing
	defer resp.Body.Close()

	// use `var` or directly `:=`
    // inline new type which is a structure
    var d struct {
        Main struct {
            Kelvin float64 `json:"temp"`    // name tag type
	                                        // tag helps to map json to this struct
        } `json:"main"`
    }

	// interface
	// json.NewDecoder takes io.Reader interface, not a concrete HTTP resp body
	// http.Response.Body satisfy io.Reader
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

    log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
    return d.Main.Kelvin, nil
}

// [!] so this is basically the cons of this kind of type system, recall me ML
type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
    sum := 0.0

    for _, provider := range w {
        k, err := provider.temperature(city)
        if err != nil {
            return 0, err
        }

        sum += k
    }

    return sum / float64(len(w)), nil
}

// http framework
func main() {
    mw := multiWeatherProvider{
        openWeatherMap{apiKey: "6c68ab3ee677cf2fe425ff3ac9ab9b6f"},
        openWeatherMap{apiKey: "354e7ad21558102181b0136a25f78c11"},
    }

	// HandlerFunc register "/" with hello for ServeMux(http router)
	// "/" matches all path, not only the "/"
	// first-class function
	// [!] So can see this function change and save some internal states we cannot see here
	http.HandleFunc("/hello", hello) // (pattern, handler that been registered)

	// at first a mistype `HandleFunc` as `HandlerFunc` and the compiler
	// report it poorly: ./main.go:62: too many arguments to conversion to http.HandlerFunc: http.HandlerFunc("/weather/", func literal)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
        begin := time.Now()
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		temp, err := mw.temperature(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
            "city": city,
            "temp": temp,
            "took": time.Since(begin).String(),
        })
	})

	http.ListenAndServe(":8080", nil)
}

// with a type system of strong, static, inferred
// [!] so why still make programmer define the type???
// http.HandlerFunc
// server will spawn a new goroutine executing this for each request
func hello(w http.ResponseWriter, r *http.Request) {
	// w.Write take general []byte, so need to cast
	// if lost ')', compiler will raise: ./main.go:19: syntax error: unexpected semicolon or newline, expecting comma or )
	// it is ok
	w.Write([]byte("hello!"))
}
