FROM golang:1.17 AS build-env
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o game

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/app/game .
ENTRYPOINT ["./game"]