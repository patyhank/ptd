package core

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hinshun/vt10x"
	"github.com/patyhank/ptd/core/event"
	"github.com/patyhank/ptd/uao"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"slices"
	"time"
)

type Client struct {
	Username, Password string
	Host, Origin       string

	registry EventRegistry
	EventManager

	conn          *websocket.Conn
	state         vt10x.Terminal
	notifyChannel map[string]chan string
	writeChan     chan byte
}

func NewConn(host, origin string) *Client {
	v := &Client{Host: host, Origin: origin}
	v.notifyChannel = make(map[string]chan string)
	v.writeChan = make(chan byte, 2048)
	v.PrepareWait("default")
	v.registry.RegisterDefaultEvents()

	return v
}

func (c *Client) Connect(ctx context.Context) {
	conn, response, _ := websocket.DefaultDialer.Dial(c.Host, http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"},
		"Origin":     []string{c.Origin},
	})
	if response.StatusCode != 101 {
		fmt.Println(response.Status)
		return
	}
	c.conn = conn
	c.state = vt10x.New()

	c.state.Resize(80, 24)

	commands := [][]byte{
		{IAC, WILL, OptTerminalType},
		slices.Concat([]byte{IAC, SB, OptTerminalType, TerminalTypeIs}, []byte("VT100"), []byte{IAC, SE}),
		{IAC, WILL, OptNegotiateAboutWindowSize},
		{IAC, SB, OptNegotiateAboutWindowSize, 0, 80, 0, 24, IAC, SE},
		{IAC, DO, OptEcho},
		{IAC, DO, OptSuppressGoAhead},
	}

	for _, command := range commands {
		err := conn.WriteMessage(websocket.BinaryMessage, command)
		if err != nil {
			log.Fatal("Error writing initial commands", err)
		}
	}

	payloads := make(chan []byte, 2048)
	go func() {
		for {
			_, npayload, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("Error reading message", err)
				return
			}
			payloads <- npayload
		}
	}()
	go func() {
		for {
			select {
			case data := <-c.writeChan:
				err := conn.WriteMessage(websocket.BinaryMessage, []byte{data})
				if err != nil {
					log.Fatal("Error writing message", err)
					return
				}
			}
		}
	}()

	var payload []byte
	doneChan := ctx.Done()
	go func() {
		for {
		LoadingContentLoop:
			for {
				select {
				case data := <-payloads:
					payload = append(payload, data...)
				case <-time.After(300 * time.Millisecond):
					if len(payload) > 0 {
						break LoadingContentLoop
					}
				case <-doneChan:
					conn.Close()
					return
				}
			}

			data := uao.Decode(payload)

			_, err := c.state.Write(data)
			if err != nil {
				fmt.Println(err)
			}

			payload = []byte{}
			screen := c.String()
			c.callScreenMatch(screen)

			os.WriteFile("screen.txt", []byte(screen), 0644)

			c.NotifyUpdated(screen)
		}
	}()
}

func (c *Client) String() string {
	return c.state.String()
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Events

func (c *Client) Screen() string {
	return c.state.String()
}

func (c *Client) PrepareWait(name string) {
	if ch, ok := c.notifyChannel[name]; ok {
		for len(ch) > 0 {
			<-ch
		}
		return
	}
	c.notifyChannel[name] = make(chan string, 1)
	return
}

func (c *Client) WaitUpdate(ctx context.Context, name string) {
	select {
	case <-ctx.Done():
		return
	case <-c.notifyChannel[name]:
		return
	}
}

func (c *Client) NotifyUpdated(screen string) {
	for _, ch := range c.notifyChannel {
		select {
		case ch <- screen:
		default:
		}
	}
}

//func (v *Client) SendAllUpdates(data string) {
//	for _, ch := range v.notifyChannel {
//		ch := ch
//		go func() {
//			for len(ch) > 0 {
//				<-ch
//			}
//
//			ch <- data
//		}()
//	}
//}
//
//func (v *Client) RegisterUpdateChan(name string) chan string {
//	if _, ok := v.notifyChannel[name]; ok {
//		return v.notifyChannel[name]
//	}
//	v.notifyChannel[name] = make(chan string, 16)
//
//	return v.notifyChannel[name]
//}
//
//func (v *Client) UnregisterUpdateChan(name string) {
//	delete(v.notifyChannel, name)
//}
//
//func (v *Client) WaitUpdatedChan(name string) string {
//	v.RegisterUpdateChan(name)
//	ch := v.notifyChannel[name]
//	for len(ch) > 0 {
//		<-ch
//	}
//	return <-ch
//}
//
//func (v *Client) ClearUpdateChan(name string) {
//	v.RegisterUpdateChan(name)
//	ch := v.notifyChannel[name]
//	for len(ch) > 0 {
//		<-ch
//	}
//}

func (c *Client) SendMessage(content string, sendReturn ...bool) {
	contents := uao.Encode(content)

	for _, t := range contents {
		c.writeChan <- t
	}
	if len(sendReturn) > 0 && sendReturn[0] {
		c.writeChan <- byte('\r')
	}
}

func (c *Client) SendReturn() {
	c.writeChan <- byte('\r')
}

func (c *Client) SendRefresh() {
	c.writeChan <- byte('\f')
}

func (c *Client) SendMultipleMessage(content []string, sendReturn ...bool) {
	for _, text := range content {
		text := uao.Encode(text)
		for _, t := range text {
			c.writeChan <- byte(t)
		}
	}
	if len(sendReturn) > 0 && sendReturn[0] {
		c.writeChan <- byte('\r')
	}
}

func (c *Client) callScreenMatch(data string) {
	for _, evt := range c.registry.MatchEvents {
		if evt.Matched(data) {
			if evt, ok := evt.(event.RegexMatchEvent); ok {
				evt.FillGroup(data)
			}
			c.DispatchEvent(c, evt)
		}
	}
}
