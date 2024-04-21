# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.20-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖项
RUN go mod download

# 复制源代码到容器中
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# 使用一个新的轻量级镜像
FROM alpine:latest  

# 安装 ca-certificates，如果您的应用需要与外部服务通信（HTTPS），这是必须的
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建者镜像中复制构建出的应用程序
COPY --from=builder /app/myapp .
COPY templates/ ./templates 
COPY uploads/ ./uploads

# 设置容器启动时执行的命令
CMD ["./myapp"]

