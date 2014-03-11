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
* curl -i 'http://changeme@localhost:3000/topics/test/0?offset=2000000&count=100'
* ab -k -p data.json -T application/json -n 100000 -c 10 -A changeme: http://localhost:3000/publish
* cat | log-shuttle -logs-url=http://changeme@localhost:3000/publish?topic=test -msgid="-"

## Related links

* https://github.com/heroku/event-manager-api
* https://github.com/heroku/log2viz
* https://github.com/heroku/log-shuttle
* https://github.com/bmizerany/lpx
* http://godoc.org/github.com/Shopify/sarama
* http://godoc.org/github.com/kr/logfmt

## TODO

* don't wait if there are no messages to consume
* http streaming publish and consume
* websockets publish and consume
* go-metrics uses mutexes on each metric, for small messages and low io/network
  latency there is more contention on those mutexes. At what point is that an issue?
* sarama refers it's logging to `Logger`. Find out where that is and how to get it going
  to stderr
* sarama supports tons of configuration options, we should support more of them
* use the json schema stuff to generate a ruby kafka http api client, :)
* Allow connection specific producer configuration via params for long
  running connections
