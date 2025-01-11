package uao

import (
	"bytes"
	"encoding/hex"
	"golang.org/x/text/encoding/unicode"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const cacheFolder = ".cache"
const pageMapping = "CodePageMapping"

const (
	mapping = "https://moztw.org/docs/big5/table/moz18-b2u.txt"
)

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
	B2U     = make(map[[2]byte]string)
	U2B     = make(map[string][2]byte)
)

func init() {
	_ = os.Mkdir(cacheFolder, 0755)

	if _, err := os.Stat(cacheFolder + "/" + pageMapping + ".txt"); os.IsNotExist(err) {
		download(mapping, pageMapping)
	}

	readMapping()
}

func download(url string, name string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(cacheFolder+"/"+name+".txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	_, err = io.Copy(file, resp.Body)
}

func readMapping() {
	file, err := os.ReadFile(cacheFolder + "/" + pageMapping + ".txt")
	if err != nil {
		return
	}
	encoding := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	decoder := encoding.NewDecoder()
	lines := strings.Split(strings.ReplaceAll(string(file), "\r\n", "\n"), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		hexBig5, _ := hex.DecodeString(parts[0][2:])
		hexUnicode, _ := hex.DecodeString(parts[1][2:])
		u, _ := decoder.Bytes(hexUnicode)
		var b [2]byte
		copy(b[:], hexBig5)

		U2B[string(u)] = b
		B2U[b] = string(u)
	}

}
func Decode(content []byte) []byte {
	s := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(s)
	s.Reset()

	var prev byte = 0xff
	for _, b := range content {
		if prev != 0xff {
			if u, ok := B2U[[2]byte{prev, b}]; ok {
				s.WriteString(u)
				prev = 0xff
				continue
			}
		}
		if b <= 0x80 {
			s.WriteByte(b)
			continue
		}
		prev = b
	}

	return s.Bytes()
}

func Encode(content string) []byte {
	s := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(s)
	s.Reset()

	for _, data := range content {
		if b, ok := U2B[string(data)]; ok {
			s.WriteByte(b[0])
			s.WriteByte(b[1])
			continue
		}
		s.WriteByte(byte(data))
	}

	return s.Bytes()
}
