package main

import (
	"acme-proxy/internal/config"
	"acme-proxy/internal/proxy"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	route := proxy.Route{Match: "/test", Target: "https://postman-echo.com/headers"}
	client := config.ProvideHTTPClient()
	http.Handle(route.Match, route.ProvideProxy(client))

	log.Println("Listening on port: ", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("listen and serve failed: %v", err)
	}
}

