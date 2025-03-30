# Stage 1: Build Go binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy source code
COPY . .

# Build binary แบบ static เพื่อให้รองรับ Alpine Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /root/

# Copy compiled binary จาก builder stage
COPY --from=builder /app/main /root/main

# ให้สิทธิ์ execute ไฟล์ main
RUN chmod +x /root/main

EXPOSE 3001

CMD ["/root/main"]
