FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ctc ./cmd/ctc/main.go

FROM alpine:latest

COPY --from=build /app/ctc /ctc

ENTRYPOINT ["/ctc"]