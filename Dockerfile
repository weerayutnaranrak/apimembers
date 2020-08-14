FROM golang:latest

LABEL maintainer="Rajeev Singh weerayut.naranrak@gmail.com"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]