package http_test

import (
	"bufio"
	"net/http"
	"strings"
	"testing"

	testdispatcher "github.com/v2ray/v2ray-core/app/dispatcher/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	. "github.com/v2ray/v2ray-core/proxy/http"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestHopByHopHeadersStrip(t *testing.T) {
	assert := assert.On(t)

	rawRequest := `GET /pkg/net/http/ HTTP/1.1
Host: golang.org
Connection: keep-alive,Foo, Bar
Foo: foo
Bar: bar
Proxy-Connection: keep-alive
Proxy-Authenticate: abc
User-Agent: Mozilla/5.0 (Macintosh; U; Intel Mac OS X; de-de) AppleWebKit/523.10.3 (KHTML, like Gecko) Version/3.0.4 Safari/523.10
Accept-Encoding: gzip
Accept-Charset: ISO-8859-1,UTF-8;q=0.7,*;q=0.7
Cache-Control: no-cache
Accept-Language: de,en;q=0.7,en-us;q=0.3

`
	b := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := http.ReadRequest(b)
	assert.Error(err).IsNil()
	assert.String(req.Header.Get("Foo")).Equals("foo")
	assert.String(req.Header.Get("Bar")).Equals("bar")
	assert.String(req.Header.Get("Connection")).Equals("keep-alive,Foo, Bar")
	assert.String(req.Header.Get("Proxy-Connection")).Equals("keep-alive")
	assert.String(req.Header.Get("Proxy-Authenticate")).Equals("abc")

	StripHopByHopHeaders(req)
	assert.String(req.Header.Get("Connection")).Equals("close")
	assert.String(req.Header.Get("Foo")).Equals("")
	assert.String(req.Header.Get("Bar")).Equals("")
	assert.String(req.Header.Get("Proxy-Connection")).Equals("")
	assert.String(req.Header.Get("Proxy-Authenticate")).Equals("")
}

func TestNormalGetRequest(t *testing.T) {
	assert := assert.On(t)

	testPacketDispatcher := testdispatcher.NewTestPacketDispatcher(nil)

	port := v2nettesting.PickPort()
	httpProxy := NewServer(&Config{}, testPacketDispatcher, v2net.LocalHostIP, port)
	defer httpProxy.Close()

	err := httpProxy.Start()
	assert.Error(err).IsNil()
	assert.Port(port).Equals(httpProxy.Port())

	httpClient := &http.Client{}
	resp, err := httpClient.Get("http://127.0.0.1:" + port.String() + "/")
	assert.Error(err).IsNil()
	assert.Int(resp.StatusCode).Equals(400)
}
