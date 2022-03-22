FROM golang:1.17

COPY . /raptor-hello

WORKDIR /raptor-hello

RUN go build

ENTRYPOINT ["./hello"]