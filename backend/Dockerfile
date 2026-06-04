FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server .

FROM alpine:3.19
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
