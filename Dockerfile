FROM alpine as god-jay-web
RUN mkdir -p /service
COPY ./web /service
COPY cmd/web/home.html /service/
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /service/
EXPOSE 81


FROM alpine as god-jay-ws
RUN mkdir -p /service
COPY ./websocket /service
COPY cmd/websocket/etc /service/etc
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /service/
EXPOSE 8081