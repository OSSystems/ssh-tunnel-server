FROM golang:alpine AS build
RUN apk add --update git && \
    go get github.com/eclipse/paho.mqtt.golang && \
    go get github.com/gliderlabs/ssh && \
    go get github.com/Sirupsen/logrus
ADD . /src
RUN cd /src && go build -o ssh-gateway

FROM alpine:3.7
EXPOSE 2222
WORKDIR /app
COPY --from=build /src/ssh-gateway /app/
ADD key.pem /app/key.pem
ENTRYPOINT ./ssh-gateway
