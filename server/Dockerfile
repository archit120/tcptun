FROM golang:1.17-bullseye

RUN apt update && apt install -y iproute2 && apt install -y iputils-ping  && \
    apt install -y iptables

WORKDIR /code

