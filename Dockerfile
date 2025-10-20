FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o main ./cmd

# ==================
# 메인 이미지
# ==================
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/GDH-Project/api"

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Seoul

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]