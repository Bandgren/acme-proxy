package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type Route struct {
	Match string
	Target  string
}

func (r *Route)ProvideProxy(client *http.Client)  http.HandlerFunc {
	upstream, err := url.Parse(r.Target)
	if err != nil {
		log.Fatal(err)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		req.Host = upstream.Host
		req.URL.Host = upstream.Host
		req.URL.Scheme = upstream.Scheme
		req.URL.Path = upstream.Path
		req.RequestURI = ""
		host, _, _ := net.SplitHostPort(req.RemoteAddr)
		req.Header.Add("X-Forwarded-For", host)
		req.Header.Add("key", "super-secret-key")

		resp, err := client.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}

		for key, values := range resp.Header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

		rw.WriteHeader(resp.StatusCode)
		_, err = io.Copy(rw, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}

