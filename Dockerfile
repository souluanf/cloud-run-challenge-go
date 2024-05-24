FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o cloud-run-challenge-go main.go

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cloud-run-challenge-go .
ENTRYPOINT ["./cloud-run-challenge-go"]
