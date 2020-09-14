package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type route struct {
	match string
	target  string
}

func NewProxy(r *route) http.Handler {
	director := func(req *http.Request) {
		out, _ := url.Parse(r.target)
		req.Header.Add("X-Forwarded-Host", out.Host)
		req.Header.Add("X-Origin-Host", out.Host)
		req.Header.Add("key", "super-secret-key")
		req.URL.Scheme = out.Scheme
		req.URL.Host = out.Host
		req.URL.Path = out.Path
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	route := route{match: "/test", target: "https://postman-echo.com/headers"}
	http.Handle(route.match, NewProxy(&route))

	log.Println("Listening on port: ", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("listen and serve failed: %v", err)
	}
}

