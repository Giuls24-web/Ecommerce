FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o floriluz .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/floriluz .
COPY --from=builder /app/frontend ./frontend
EXPOSE 8080
CMD ["./floriluz"]