all: build-go build-dockerfile docker-compose

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
	docker compose -p god-jay-chat up -d reverse-proxy
	docker compose -p god-jay-chat up -d web
	docker compose -p god-jay-chat up -d --scale ws=3


### destroy
stop-docker-compose:
	docker compose -p god-jay-chat down

rm-build-go:
	rm -rf web websocket

rm-build-dockerfile:
	docker rmi god-jay-web || true
	docker rmi god-jay-ws || true