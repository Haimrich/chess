FROM golang:alpine as build-env
WORKDIR /go/src/app
ENV GOOS=linux GOARCH=amd64 
COPY . .
RUN apk update && apk add gcc librdkafka-dev openssl-libs-static zlib-static zstd-libs libsasl librdkafka-static lz4-dev lz4-static zstd-static libc-dev musl-dev 
RUN go build -tags musl -ldflags '-w -extldflags "-static"' -o challenge

#FROM alpine
#WORKDIR /app
#COPY --from=build-env /go/src/app/wsnode .
EXPOSE 8081
ENTRYPOINT ["./challenge"]
