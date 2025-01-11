package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hinshun/vt10x"
	"github.com/patyhank/ptd/uao"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	// Telnet commands. rfc854, rfc1116, rfc1123
	EOF   = 236 // end of file
	SUSP  = 237 // suspend process
	ABORT = 238 // abort process
	EOR   = 239 // end of record (transparent mode, used for prompt marking)
	SE    = 240 // end sub negotiation
	NOP   = 241 // nop (used for keep alive messages	)
	DM    = 242 // data mark--for connect. cleaning
	BREAK = 243 // break
	IP    = 244 // interrupt process (permanently)
	AO    = 245 // abort output (but let program finish)
	AYT   = 246 // are you there
	EC    = 247 // erase the current character
	EL    = 248 // erase the current line
	GA    = 249 // you may reverse the line (used for prompt marking)
	SB    = 250 // interpret as subnegotiation
	WILL  = 251 // I will use option
	WONT  = 252 // I won"t use option
	DO    = 253 // please, you use option
	DONT  = 254 // you are not to use option
	IAC   = 255 // interpret as command

	// Telnet Options. rfc855
	OptTM         = 6   // timing mark. rfc860
	OptTType      = 24  // terminal type. rfc930, rfc1091
	OptEOR        = 25  // end of record. rfc885
	OptNAWS       = 31  // negotiate about window size. rfc1073
	OptLineMode   = 34  // linemode. rfc1184
	OptEnviron    = 36  // environment option. rfc1408
	OptNewEnviron = 39  // new environment option. rfc1572
	OptCharset    = 42  // character set. rfc2066
	OptMSDP       = 69  // mud server data protocol. @see: https://tintin.sourceforge.io/protocols/msdp/
	OptMSSP       = 70  // mud server status protocol. @see: https://tintin.sourceforge.io/protocols/mssp/
	OptMCCP       = 86  // mud client compression protocol(v2). @see: https://tintin.sourceforge.io/protocols/mccp/
	OptMSP        = 90  // mud sound protocol. @see: https://www.zuggsoft.com/zmud/msp.htm
	OptMXP        = 91  // mud extension protocol. @see: https://www.zuggsoft.com/zmud/mxp.htm
	OptATCP       = 200 // achaea telnet client protocol. @see: https://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html
	OptGMCP       = 201 // generic mud client protocol. @see: https://tintin.sourceforge.io/protocols/gmcp/

	// OptTType
	TTypeIs   = 0
	TTypeSend = 1

	// MTTS standard codes @see: https://tintin.sourceforge.io/protocols/mtts/
	TTypeANSI            = 1
	TTypeVT100           = 2
	TTypeUTF8            = 4
	TType256Colors       = 8
	TTypeMouseTracking   = 16
	TTypeOscColorPalette = 32
	TTypeScreenReader    = 64
	TTypeProxy           = 128
	TTypeTrueColor       = 256
	TTypeMNES            = 512
	TTypeMSLP            = 1024

	// OptEnviron, OptNewEnviron
	EnvironIs      = 0
	EnvironSend    = 1
	EnvironVar     = 0
	EnvironValue   = 1
	EnvironESC     = 2
	EnvironUserVar = 3
)

type DetectingTarget struct {
	Messages     [][]string
	Regexes      [][]*regexp.Regexp
	JustOneRegex bool // whether to match just one regex then call the function and break
}

type MessageHandler struct {
	Detect     *DetectingTarget
	Priority   int
	NoGo       bool                 // whether to run the function in a goroutine
	Call       func() error         // function to call
	RegexCall  func([]string) error // function to call with regex matches
	ScreenCall func(string) error   // function to call with screen content
}

type Conn struct {
	Conn                 *websocket.Conn
	State                vt10x.Terminal
	registeredUpdateChan map[string]chan string
	Handlers             []MessageHandler
	writeChan            chan byte
}

func (v *Conn) PreInit() {
	v.registeredUpdateChan = make(map[string]chan string)
	v.writeChan = make(chan byte, 2048)
	v.RegisterUpdateChan("default")
}

func (v *Conn) Connect(host, origin string) {
	conn, response, _ := websocket.DefaultDialer.Dial(host, http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"},
		"Origin":     []string{origin},
	})
	if response.StatusCode != 101 {
		fmt.Println(response.Status)
		return
	}
	v.Conn = conn
	v.State = vt10x.New()
	go func() {
		commands := [][]byte{
			[]byte{IAC, WILL, OptTType},
			slices.Concat([]byte{IAC, SB, OptTType, TTypeIs}, []byte("VT100"), []byte{IAC, SE}),
			[]byte{IAC, WILL, OptNAWS},
			//[]byte{IAC, SB, OptNAWS, 0, 200, 0, 200, IAC, SE},
			[]byte{IAC, DO, 1},
			[]byte{IAC, DO, 3},
			[]byte{IAC, WONT, 0},
			[]byte{IAC, DONT, 0},
		}
		for _, command := range commands {
			conn.WriteMessage(websocket.BinaryMessage, command)
		}
	}()
	v.init(conn)
}

func (v *Conn) init(conn *websocket.Conn) {
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
			case data := <-v.writeChan:
				err := conn.WriteMessage(websocket.BinaryMessage, []byte{data})
				if err != nil {
					log.Fatal("Error writing message", err)
					return
				}
			}
		}
	}()

	v.State.Resize(80, 24)

	var payload []byte
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
			}
		}

		data := uao.Decode(payload)

		_, err := v.State.Write(data)
		if err != nil {
			fmt.Println(err)
		}

		payload = []byte{}
		screen := v.String()
		//fmt.Println(screen)
		for _, handler := range v.Handlers {
			for _, messages := range handler.Detect.Messages {
				if containsAllConditions(screen, messages) {
					handler.callData(screen, nil)
					break
				}
			}

			for _, regexes := range handler.Detect.Regexes {
				for _, regex := range regexes {
					if regex.MatchString(screen) {
						matches := regex.FindAllStringSubmatch(screen, -1)
						for _, match := range matches {
							handler.callData(screen, match)
						}
						if handler.Detect.JustOneRegex {
							break
						}
					}
				}
			}
		}

		v.SendAllUpdates(screen)
	}
}

func (h *MessageHandler) callData(screen string, match []string) {
	if !h.NoGo {
		if h.RegexCall != nil && match != nil {
			go h.RegexCall(match)
		}

		if h.Call != nil {
			go h.Call()
		}
		if h.ScreenCall != nil {
			go h.ScreenCall(screen)
		}
	} else {
		if h.RegexCall != nil && match != nil {
			h.RegexCall(match)
		}
		if h.Call != nil {
			h.Call()
		}
		if h.ScreenCall != nil {
			h.ScreenCall(screen)
		}
	}

}
func (v *Conn) String() string {
	return v.State.String()
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func unwrap[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func (v *Conn) AddListener(listeners ...MessageHandler) {
	v.Handlers = append(v.Handlers, listeners...)
	sort.SliceStable(v.Handlers, func(i, j int) bool {
		return v.Handlers[i].Priority > v.Handlers[j].Priority
	})
}

func (v *Conn) SendAllUpdates(data string) {
	for _, ch := range v.registeredUpdateChan {
		ch := ch
		go func() {
			for len(ch) > 0 {
				<-ch
			}

			ch <- data
		}()
	}
}

func (v *Conn) RegisterUpdateChan(name string) chan string {
	if _, ok := v.registeredUpdateChan[name]; ok {
		return v.registeredUpdateChan[name]
	}
	v.registeredUpdateChan[name] = make(chan string, 16)

	return v.registeredUpdateChan[name]
}

func (v *Conn) UnregisterUpdateChan(name string) {
	delete(v.registeredUpdateChan, name)
}

func (v *Conn) WaitUpdatedChan(name string) string {
	v.RegisterUpdateChan(name)
	ch := v.registeredUpdateChan[name]
	for len(ch) > 0 {
		<-ch
	}
	return <-ch
}

func (v *Conn) ClearUpdateChan(name string) {
	v.RegisterUpdateChan(name)
	ch := v.registeredUpdateChan[name]
	for len(ch) > 0 {
		<-ch
	}
}

func (v *Conn) SendMessage(content string, newline bool) {
	for _, c := range content {
		v.writeChan <- byte(c)
	}
	if newline {
		v.writeChan <- byte('\r')
	}
}

func (v *Conn) SendMultipleMessage(content []string, newline bool) {
	for _, text := range content {
		for _, c := range text {
			v.writeChan <- byte(c)
		}
	}
	if newline {
		v.writeChan <- byte('\r')
	}
}

func containsAllConditions(content string, cond []string) bool {
	for _, s := range cond {
		if !strings.Contains(content, s) {
			return false
		}
	}
	return true
}
