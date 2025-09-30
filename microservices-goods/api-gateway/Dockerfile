# 使用轻量级的 alpine 镜像
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# 复制预构建的二进制文件
COPY gateway-api .

# 复制配置文件
COPY etc ./etc

# 更改文件所有者
RUN chown -R appuser:appgroup /root/ && \
    chmod +x /root/gateway-api

# 切换到非root用户
USER appuser

# 暴露端口（对外统一入口）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/goods/room?room_id=1 || exit 1

# 启动应用
CMD ["./gateway-api", "-f", "etc/gateway-api.yaml"]
