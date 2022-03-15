FROM golang:1.15

COPY . /raptor-hello

WORKDIR /raptor-hello

RUN go build

ENTRYPOINT ["./hello"]