package proxy

import (
	"acme-proxy/internal/config"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type PostmanResp struct {
	Headers struct {
		XForwardedProto string `json:"x-forwarded-proto"`
		Key             string `json:"key"`
	}
}

func TestProxyHandler(t *testing.T) {
	assert := assert.New(t)
	route := Route{Match: "/test", Target: "https://postman-echo.com/headers"}
	client := config.ProvideHTTPClient()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, route.Match, nil)

	p := route.ProvideProxy(client)
	p.ServeHTTP(w, r)

	assert.Equal(w.Result().StatusCode, http.StatusOK)
	bodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	resp := PostmanResp{}
	err := json.Unmarshal(bodyBytes, &resp)
	assert.NoError(err)
	assert.Equal(resp.Headers.Key, "super-secret-key")
	assert.Equal(resp.Headers.XForwardedProto, "https")
}

func TestProxyClientTimeout(t *testing.T) {
	assert := assert.New(t)
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	route := Route{Match: "/local", Target: server.URL}
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, route.Match, nil)

	p := route.ProvideProxy(client)
	p.ServeHTTP(w, r)
	assert.Equal(w.Result().StatusCode, http.StatusInternalServerError)
	bodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	assert.True(strings.Contains(string(bodyBytes), "context deadline exceeded (Client.Timeout exceeded while awaiting headers)"))
}

func BenchmarkProxy(b *testing.B) {

	route := Route{Match: "/test", Target: "https://postman-echo.com/headers"}
	client := config.ProvideHTTPClient()

	p := route.ProvideProxy(client)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, route.Match, nil)
			p.ServeHTTP(w, r)
		}
	})

}
