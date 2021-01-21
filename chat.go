package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"strings"
	"time"
	"bnsvc.net/tmi/ircon"
)

type handler struct {
	Con *ircon.IRCon
	Channels []string
}

func (h *handler) Connected() {
	fmt.Println("#", time.Now(), "Connected")
	go func() {
		N := 50
		for i := 0; i < len(h.Channels); i += N {
			end := i + N
			if end > len(h.Channels) {
				end = len(h.Channels)
			}
			clist := strings.Join(h.Channels[i:end], ",")
			fmt.Println("# Join:", clist)
			h.Con.Send("JOIN " + clist)
			time.Sleep(time.Second * 16)
		}
	}()
}
func (handler) Disconnected(err error) {
	fmt.Println("#", time.Now(), "Disconnected:", err)
}

func (handler) Message(msg *ircon.Message) {
//	fmt.Println(msg.Raw())
	fmt.Println(msg.Args)
}

func main() {
	username := os.Getenv("TWITCH_USERNAME")
	passwd := os.Getenv("TWITCH_OAUTH_TOKEN")
	channel := "neutraldread"

	ctx := context.Background()
	con := ircon.New(username, "oauth:"+passwd)
	h := &handler {
		Con: con,
	}
	if channel != "" {
		chans := strings.Split(channel, ",")
		for i, cname := range chans {
			chans[i] = addPrefix(cname, "#")
		}
		h.Channels = chans
	}
	con.Handler = h
	con.Background(ctx)

	raw := func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Raw string `json:"raw"`
		}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&req); err != nil || req.Raw == "" {
			w.WriteHeader(400)
			return
		}
		con.Send(req.Raw)
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/raw", raw)
		s := &http.Server{
			Addr:           "localhost:2048",
			Handler:        mux,
			ReadTimeout:    120 * time.Second,
			WriteTimeout:   120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		s.ListenAndServe()
	}()

	<-ctx.Done()
}

func addPrefix(s, pfx string) string {
	if !strings.HasPrefix(s, pfx) {
		return pfx + s
	}
	return s
}
