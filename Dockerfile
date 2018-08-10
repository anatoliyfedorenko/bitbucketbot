FROM golang:1.7

RUN mkdir -p /bot

WORKDIR /bot

ADD . /bot

RUN go build .

CMD ["./bot"]