# 第一阶段：构建 Go 应用
FROM golang:1.24-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o play-tools .

# 第二阶段：运行环境
FROM alpine:latest

# 设置时区为中国
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制二进制文件和配置
COPY --from=builder /app/play-tools .
COPY --from=builder /app/config.yaml .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./play-tools"]
