all: clean-all build-go build-dockerfile docker-compose

clean-all: stop-docker-compose rm-build-go rm-build-dockerfile

clean-docker: stop-docker-compose rm-build-dockerfile

### build
build-go:
	go build -o web cmd/web/*.go
	go build -o websocket cmd/websocket/*.go

build-dockerfile:
	docker build -t god-jay-web --target god-jay-web .
	docker build -t god-jay-ws --target god-jay-ws .

docker-compose:
	#docker compose -p god-jay-chat -f docker-compose.yml -f docker-compose.kafka.yml up -d --scale ws=3
	docker compose -p god-jay-chat -f docker-compose.yml -f docker-compose.nats.yml -f docker-compose.prom-graf.yml up -d --scale ws=3


### test api
run-k6-direct:
	docker run --add-host=ws.localhost:172.17.0.1 --network god-jay-chat-net --rm -i grafana/k6 run --vus 1000 --duration 10s - <test/script-direct.js
	#docker run --add-host=ws.localhost:192.168.2.113 --network god-jay-chat-net --rm -i grafana/k6 run --vus 200 --duration 10s - <test/script.js

	#docker run --network host --rm -i grafana/k6 run --vus 200 --duration 10s - <test/script.js
	#docker run --network god-jay-chat-net --rm -i grafana/k6 run - <test/script.js

run-k6-traefik:
	docker run --add-host=ws.localhost:172.17.0.1 --network god-jay-chat-net --rm -i grafana/k6 run --vus 1000 --duration 10s - <test/script-traefik.js

run-k6-nginx:
	docker run --add-host=ws.localhost:172.17.0.1 --network god-jay-chat-net --rm -i grafana/k6 run --vus 1000 --duration 10s - <test/script-nginx.js


### destroy
stop-docker-compose:
	docker compose -p god-jay-chat down

rm-build-go:
	rm -rf web websocket

rm-build-dockerfile:
	docker rmi god-jay-web || true
	docker rmi god-jay-ws || true

prometheus-grafana:
	cd ./docker/prometheus-grafana && make