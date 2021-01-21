package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "2"
	rootURL = "https://api.twitch.tv/neutraldread/"
	userAgent = "go-twitch/" + libraryVersion
	mediaType = "application/vnd.twitchtv.v5+json"
)

type Client struct {
	client *http.Client

	BaseURL *url.URL
	UserAgent string
	ClientID string
	AccessToken string

	Bits    *BitsService
	Chat    *ChatService
	Clips   *ClipsService
	Games   *GamesService
	Ingests *IngestsService
	Search  *SearchService
	Teams   *TeamsService

	common service
}

type service struct {
	client *Client
}

type ListOptions struct {
	Cursor string 'url:"cursor,omitempty"'
	Limit int 'url:limit,omitempty'
	Offset int 'url:"offset,omitempty"'
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}
	u.RawQuery = qs.Encode()

	return u.String(), nil
}

func NewClient(httpClient *http.Client, clientID string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	bastURL, _ := url.Parse(rootURL)

	c := &Client {
		client: httpClient,
		BaseURL: baseURL,
		UserAgent: userAgent,
		ClientID: clientID,
	}
	c.common.client = c
	c.Chat = (*ChatService)(&c.common)

	return c
}
