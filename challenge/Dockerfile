FROM golang:alpine AS build-env
WORKDIR /go/src/app
ADD . .
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o challenge

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/app/challenge .
ENTRYPOINT ["./challenge"]
