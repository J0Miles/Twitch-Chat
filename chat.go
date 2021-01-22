package main

import (
	"strconv"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"strings"
	"time"
	"bnsvc.net/tmi/ircon"
	"github.com/gookit/color"
)

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type handler struct {
	Con *ircon.IRCon
	Channels []string
}

type MsgObj struct {
	Color string `json:"color"`
	RoomID string `json:"room-id"`
	Subscriber string `json:"subscriber"`
}

func (h *handler) Connected() {
	fmt.Println("->", time.Now(), "Connected")
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


func ParseHex(class string) (color.RGBColor) {
	var rgb RGB                              
	var err error
	s := strings.Split(class, ";")
	for i := range s {
		cs := strings.Contains(s[i], "color=")
		if cs {
		 hexValue := strings.Split(s[i], "=")
		 var removeHash = strings.Split(hexValue[1], "#")
		 justHex := strings.Join(removeHash, "")
		 var hex Hex = Hex(justHex)
		 rgb, err = Hex2RGB(hex)
		 if err != nil {
		panic("Couldn't convert hex to rgb")
	}
	 }
}
c := color.RGB(rgb.Red, rgb.Green, rgb.Blue)
return c
}

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)
	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}

func (handler) Message(msg *ircon.Message) {
	rgbColor := ParseHex(msg.Raw())
	rgbColor.Println(msg.Args)
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
