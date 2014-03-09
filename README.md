Don't use this for anything. It's just my play ground for go and kafka.

## Useful commands

* rake build
* rake topic:create
* foreman start -f Procfile.dev
* rake agent:consumer
* curl -i -d '{"topic": "test", "messages": [ "test", "test2" ]}' http://localhost:9393/publish
* curl -i http://localhost:9393/topics/test
* curl -i http://localhost:9393/topics/test/0
* ab -k -p data.json -n 100000 -c 10 http://localhost:3000/publish
* cat | log-shuttle -logs-url=http://localhost:9393/syslog -msgid="-"

## Related links

* https://github.com/heroku/event-manager-api
* https://github.com/heroku/log2viz
* https://github.com/heroku/log-shuttle
* https://github.com/bmizerany/lpx
* http://godoc.org/github.com/Shopify/sarama
* http://godoc.org/github.com/kr/logfmt

## TODO

* configuration
* async threading, there will be lots of long running open connections
* authentication
* better creation of topics
* streaming HTTP
* websockets
* topic name patterns for consuming and multiple topics per publish
* internal metrics

