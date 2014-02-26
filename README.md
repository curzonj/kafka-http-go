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
