FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /light-api cmd/server/main.go

EXPOSE 8080

ENTRYPOINT [ "/light-api" ]

