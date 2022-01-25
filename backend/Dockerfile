#FROM golang:alpine AS build-env
#RUN apk --no-cache add build-base git bzr mercurial gcc
#ADD . /src
#RUN cd /src && go build -o backend

FROM golang:alpine AS build-env
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend

FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/app/backend .
EXPOSE 8080
ENTRYPOINT ./backend
