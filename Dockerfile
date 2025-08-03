FROM golang:1.21.6

WORKDIR /app

COPY ./src /app

RUN go mod download

CMD ["go", "run", "main.go"]
