#---- Libreria Pistache

FROM ubuntu:focal as pistache-env

ENV DEBIAN_FRONTEND=noninteractive 
RUN apt-get update && apt-get install -y software-properties-common
RUN add-apt-repository ppa:pistache+team/unstable
RUN apt update
RUN apt install -y libpistache-dev

#---- Compilazione

FROM pistache-env as build-env
RUN apt-get install -y gcc g++ make cmake pkg-config

WORKDIR /app/
COPY . /app
RUN mkdir -p build && cd build && rm -r * && cmake .. && make

#---- Esecuzione

FROM pistache-env as final
WORKDIR /app
COPY --from=build-env /app/build/engine .
EXPOSE 9080
ENTRYPOINT ["./engine"]

# serve libreria quindi niente alpine
#FROM alpine
#WORKDIR /app
#COPY --from=build-env /app/build/engine .
#EXPOSE 8090
#ENTRYPOINT ./engine
