package source

import (
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

const preamble = "\x00\x00\x00\x7b\x01\x00\x00\x00\x01\x00\x00\x12\x73\x75\x70\x70" +
	"\x6f\x72\x74\x65\x64\x2d\x76\x65\x72\x73\x69\x6f\x6e\x73\x08\x03" +
	"\x32\x2e\x30\x0e\x6d\x61\x78\x2d\x66\x72\x61\x6d\x65\x2d\x73\x69" +
	"\x7a\x65\x03\xfc\xf0\x06\x0c\x63\x61\x70\x61\x62\x69\x6c\x69\x74" +
	"\x69\x65\x73\x08\x0a\x70\x69\x70\x65\x6c\x69\x6e\x69\x6e\x67\x09" +
	"\x65\x6e\x67\x69\x6e\x65\x2d\x69\x64\x08\x24\x38\x41\x38\x33\x35" +
	"\x32\x39\x45\x2d\x42\x33\x34\x30\x2d\x34\x42\x38\x37\x2d\x39\x34" +
	"\x41\x45\x2d\x42\x44\x38\x31\x36\x46\x42\x45\x39\x33\x32\x46"

type spoeTestCase struct {
	name     string
	spoe     []byte
	expected mirror.Request
}

var spoeTestCases = []spoeTestCase{
	{
		name: "GET request",
		spoe: []byte("\x00\x00\x00\x8d\x03\x00\x00\x00\x01\x07\x01\x06\x6d\x69\x72\x72" +
			"\x6f\x72\x05\x06\x6d\x65\x74\x68\x6f\x64\x08\x03\x47\x45\x54\x04" +
			"\x70\x61\x74\x68\x08\x09\x2f\x74\x68\x65\x2f\x70\x61\x74\x68\x03" +
			"\x76\x65\x72\x08\x03\x31\x2e\x31\x07\x68\x65\x61\x64\x65\x72\x73" +
			"\x09\x48\x04\x48\x6f\x73\x74\x0f\x31\x32\x37\x2e\x30\x2e\x30\x2e" +
			"\x31\x3a\x31\x30\x30\x38\x30\x0a\x55\x73\x65\x72\x2d\x41\x67\x65" +
			"\x6e\x74\x0b\x63\x75\x72\x6c\x2f\x37\x2e\x36\x34\x2e\x30\x06\x41" +
			"\x63\x63\x65\x70\x74\x03\x2a\x2f\x2a\x08\x58\x2d\x48\x65\x61\x64" +
			"\x65\x72\x05\x76\x61\x6c\x75\x65\x00\x00\x04\x62\x6f\x64\x79\x09\x00"),
		expected: mirror.Request{
			Method:      mirror.Method_GET,
			Path:        "/the/path",
			HttpVersion: mirror.HTTPVersion_HTTP1_1,
			Headers: map[string]*mirror.HeaderValue{
				"Host":       {Values: []string{"127.0.0.1:10080"}},
				"User-Agent": {Values: []string{"curl/7.64.0"}},
				"Accept":     {Values: []string{"*/*"}},
				"X-Header":   {Values: []string{"value"}},
			},
			Body: []byte{},
		},
	},
	{
		name: "POST request",
		spoe: []byte("\x00\x00\x00\xba\x03\x00\x00\x00\x01\x05\x01\x06\x6d\x69\x72\x72" +
			"\x6f\x72\x05\x06\x6d\x65\x74\x68\x6f\x64\x08\x04\x50\x4f\x53\x54" +
			"\x04\x70\x61\x74\x68\x08\x01\x2f\x03\x76\x65\x72\x08\x03\x31\x2e" +
			"\x31\x07\x68\x65\x61\x64\x65\x72\x73\x09\x79\x04\x48\x6f\x73\x74" +
			"\x0f\x31\x32\x37\x2e\x30\x2e\x30\x2e\x31\x3a\x31\x30\x30\x38\x30" +
			"\x0a\x55\x73\x65\x72\x2d\x41\x67\x65\x6e\x74\x0b\x63\x75\x72\x6c" +
			"\x2f\x37\x2e\x36\x34\x2e\x30\x06\x41\x63\x63\x65\x70\x74\x03\x2a" +
			"\x2f\x2a\x0e\x43\x6f\x6e\x74\x65\x6e\x74\x2d\x4c\x65\x6e\x67\x74" +
			"\x68\x01\x33\x0c\x43\x6f\x6e\x74\x65\x6e\x74\x2d\x54\x79\x70\x65" +
			"\x21\x61\x70\x70\x6c\x69\x63\x61\x74\x69\x6f\x6e\x2f\x78\x2d\x77" +
			"\x77\x77\x2d\x66\x6f\x72\x6d\x2d\x75\x72\x6c\x65\x6e\x63\x6f\x64" +
			"\x65\x64\x00\x00\x04\x62\x6f\x64\x79\x09\x03\x48\x45\x59"),
		expected: mirror.Request{
			Method:      mirror.Method_POST,
			Path:        "/",
			HttpVersion: mirror.HTTPVersion_HTTP1_1,
			Headers: map[string]*mirror.HeaderValue{
				"Host":           {Values: []string{"127.0.0.1:10080"}},
				"User-Agent":     {Values: []string{"curl/7.64.0"}},
				"Accept":         {Values: []string{"*/*"}},
				"Content-Length": {Values: []string{"3"}},
				"Content-Type":   {Values: []string{"application/x-www-form-urlencoded"}},
			},
			Body: []byte("HEY"),
		},
	},
}

func TestHAProxySPOE(t *testing.T) {
	for _, testCase := range spoeTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			name, err := ioutil.TempDir("/tmp", "http-mirror-test")
			require.NoError(t, err)
			defer os.RemoveAll(name)

			mod, err := NewHAProxySPOE(mirror.ModuleContext{}, []byte(`{"listen_addr": "@`+name+`spoe.sock"}`))
			require.NoError(t, err)

			go func() {
				conn, err := net.Dial("unix", name+"spoe.sock")
				require.NoError(t, err)

				conn.Write([]byte(preamble))
				conn.Write(testCase.spoe)
			}()

			req := <-mod.Output()
			require.Equal(t, testCase.expected, req)
		})
	}
}
