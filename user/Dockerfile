FROM golang:alpine AS build-env
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/app/user .
EXPOSE 8080
EXPOSE 8070
EXPOSE 2112
ENTRYPOINT ./user
