# 構建階段
FROM golang:1.23.5-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 安裝必要的構建工具
RUN apk add --no-cache git

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 運行階段
FROM alpine:latest

# 安裝必要的運行時依賴
RUN apk add --no-cache ca-certificates tzdata

# 設置工作目錄
WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# 創建必要的目錄
RUN mkdir -p dist

# 設置時區
ENV TZ=Asia/Taipei

# 暴露端口
EXPOSE 61018

# 運行應用程序
CMD ["./main"] 