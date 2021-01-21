package main

import (
	"fmt"
	"net/http"
	"net/url"
	// "log"
	// "html"
	// "os"
	"encoding/json"
	"io/ioutil"
)

func main() {
// go func() {
//		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
//		})
//		log.Println("Serving...")
//		log.Fatal(http.ListenAndServe(":8080", nil))
//	}()

	v := url.Values{}
	v.Set("client_id", "kxt99cdj982ul02ldwgblq87qk832m")
	v.Set("redirect_uri", "http://localhost:8080/")
	v.Set("response_type", "token")
	v.Set("scope", "user_read chat:read chat:edit channel:moderate channel_read")
	fmt.Println("https://id.twitch.tv/oauth2/authorize?"+v.Encode())
//	resp, err := http.Get("https://id.twitch.tv/oauth2/authorize?"+v.Encode())
//	if err != nil {
//		panic(err)
//	}
//	err = checkResponse(resp)
//	if err != nil {
//		panic(err)
//	}
// var r map[string]interface{}
//	 err = json.NewDecoder(resp.Body).Decode(&r)
//	 if err != nil {
//		 panic(err)
//	 }
//	 fmt.Println(resp)
}

type ErrorResponse struct {
	// HTTP response that cause this error.
	Response *http.Response

	// Error message.
	Message string `json:"message,omitempty"`
}

func checkResponse(r *http.Response) error {
	if 200 <= r.StatusCode && r.StatusCode <= 299 {
		fmt.Println("Success")
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func (e *ErrorResponse) Error() string {
	r := e.Response

	return fmt.Sprintf("%v %v: %d %v",
		r.Request.Method, r.Request.URL, r.StatusCode, e.Message)
}

