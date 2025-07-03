FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build ./cmd/bridget/.

FROM alpine AS deploy

WORKDIR /app

COPY --from=builder /app/bridget/. . 

ENTRYPOINT [ "./bridget" ]
