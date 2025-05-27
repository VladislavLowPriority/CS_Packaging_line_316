FROM golang:1.24.1

WORKDIR /app

COPY ["./", "./"]

RUN go build -o bin/manageSys ./cmd/main.go

ENTRYPOINT [ "/app/bin/manageSys" ]