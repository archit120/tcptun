FROM golang:1.17-bullseye

RUN apt update && apt install iproute2
WORKDIR /code

