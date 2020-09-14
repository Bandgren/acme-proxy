# acme-proxy

**THE BEER-WARE LICENSE**

Bandgren wrote this file. As long as you retain this notice you can do whatever you want with this stuff. If we meet some day, and you think this stuff is worth it, you can buy me a beer in return.

License originally authored by Poul-Henning Kamp (phk).

## Run acme-proxy
``` 
go run main.go -p X
```
Where `X` defaults to port 8080.

## Test it out

The `/test` endpoint proxies the request towards `https://postman-echo.com/headers`
```
curl http://localhost:8080/test 
```
And you will get the following JSON response
```
{
    "headers":
        {
            "x-forwarded-proto":"https",
            "x-forwarded-port":"443",
            "host":"localhost",
            "x-amzn-trace-id":"Root=1-5f5f2cdc-6077a6677943234e21c89ea2",
            "user-agent":"curl/7.72.0",
            "accept":"*/*",
            "x-forwarded-host":"",
            "x-origin-host":"postman-echo.com",
            "key":"super-secret-key",
            "accept-encoding":"gzip"
        }
}
```

## Things to think about:
* How to protect the key?
    * Depending on the domain, if running inside kubernetes, as a kube secret, otherwise a HMAC or OAuth solution.
* Limitations
    * Currently, only supports one endpoint to proxy match.
    * Using standard go std lib reverseProxy. Maybe not super flexible for future special business cases.
* Possible tests
    * Since the only "business" logic in this case is to append the `key` header this could be tested for.
    But since we are using well maintained std lib code that is thoroughly tested adding more "same" tests Would be abit redundant. 
    
