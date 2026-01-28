FROM golang:1.24 AS builder
WORKDIR /usr/src/app
COPY . .
WORKDIR /usr/src/app/microservices/shipping
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipping ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /usr/src/app/microservices/shipping/shipping .
EXPOSE 3002
CMD ["./shipping"]
