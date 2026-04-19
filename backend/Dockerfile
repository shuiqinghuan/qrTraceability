FROM golang:1.26-alpine as builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8000

CMD ["./api"]
