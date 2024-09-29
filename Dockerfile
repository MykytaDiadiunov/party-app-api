FROM golang:1.23.1-alpine

RUN apk add --no-cache git

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY . .

RUN go mod download

RUN mkdir /main
RUN go build -o /main/main ./cmd/server/main.go

CMD ["/main/main"]