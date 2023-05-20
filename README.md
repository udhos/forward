# forward

```
# start miniapi on port 2000
export ADDR=:2000
miniapi
```

```
# start forward
forward
```

```
# call forward

curl -d '{"url":"http://localhost:2000/v1/hello","method":"PUT","set_headers":{"a":"b"},"body":"aaaaa"}' localhost:8080/forward

{"request":{"headers":{"A":["b"],"Accept-Encoding":["gzip"],"Content-Length":["5"],"User-Agent":["Go-http-client/1.1"],"X-B3-Sampled":["1"],"X-B3-Spanid":["2e0edda835b4b13f"],"X-B3-Traceid":["300586204b1d49a4667952f629fbc748"]},"method":"PUT","uri":"/v1/hello","host":"localhost:2000","body":"aaaaa","form_query":{},"form_post":{},"parameters":{"param1":"","param2":""}},"message":"ok","status":200,"server_hostname":"ubuntu","server_version":"1.0.5"}
```
