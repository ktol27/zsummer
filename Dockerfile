# 构建阶段
FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
RUN mkdir -p /app/avatars
# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=build /app/main .
COPY init.sql .

# 设置环境变量
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=982655
ENV DB_NAME=02
ENV CHOKIDAR_USEPOLLING=true

# 将头像目录映射到宿主机，以确保头像存储持久化
VOLUME ["/app/avatars"]

EXPOSE 8080
CMD ["./main"]

