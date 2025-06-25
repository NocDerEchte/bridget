FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \ 
    go build ./cmd/bridget/.


FROM alpine

WORKDIR /app

COPY --from=builder /app/bridget/. . 

ENTRYPOINT [ "./bridget" ]
