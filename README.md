Don't use this for anything. It's just my play ground for go and kafka.

## env var defaults

```
BASIC_AUTH=changeme
PORT=3000
```

## Useful commands

* rake build
* foreman start -f Procfile.dev
* rake agent:consumer
* go run *.go
* curl -i -H "Content-Type: application/json" -d '{"topic": "test", "messages": [ "test", "test2" ]}' http://changeme@localhost:3000/publish
* curl -i http://changeme@localhost:3000/topics/test
* curl -i http://changeme@localhost:3000/topics/test/0
* ab -k -p data.json -T application/json -n 100000 -c 10 -A changeme: http://localhost:3000/publish
* cat | log-shuttle -logs-url=http://changeme@localhost:3000/publish -msgid="-"

## Related links

* https://github.com/heroku/event-manager-api
* https://github.com/heroku/log2viz
* https://github.com/heroku/log-shuttle
* https://github.com/bmizerany/lpx
* http://godoc.org/github.com/Shopify/sarama
* http://godoc.org/github.com/kr/logfmt

## TODO

* configuration
* authentication
* internal metrics
* http polling topics
* do we need to support more that 1 connection per goroutine?

