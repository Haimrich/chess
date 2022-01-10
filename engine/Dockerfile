FROM ubuntu:focal

WORKDIR /

ENV DEBIAN_FRONTEND=noninteractive 
RUN apt-get update && apt-get install -y software-properties-common gcc g++ make cmake pkg-config
RUN add-apt-repository ppa:pistache+team/unstable
RUN apt update
RUN apt install -y libpistache-dev

LABEL Name=engine Version=0.0.1
