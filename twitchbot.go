package main

import (
	"os"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"fmt"
)

const ChannelID = "76365342"

func main() {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/channels/"+ChannelID, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("User-Agent", "neutraldread")
	req.Header.Set("Client-ID", "kxt99cdj982ul02ldwgblq87qk832m")
	req.Header.Set("Authorization", "OAuth "+os.Getenv("TWITCH_OAUTH_TOKEN"))

	var resp map[string]interface{}
	_, err = Do(req, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func Do(req *http.Request, r interface{}) (*http.Response, error) {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkResponse(resp); err != nil {
		return resp, err
	}

	if r != nil {
		err = json.NewDecoder(resp.Body).Decode(r)
		if err == io.EOF {
			err = nil
		}
	}
	return resp, err
}

	type ErrorResponse struct {
	// HTTP response that cause this error.
	Response *http.Response

	// Error message.
	Message string `json:"message,omitempty"`
}

func checkResponse(r *http.Response) error {
	if 200 <= r.StatusCode && r.StatusCode <= 299 {
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

