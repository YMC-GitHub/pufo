#==============================================
# 官方 Golang 1.26.1 镜像（固定版本）
#==============================================
FROM golang:1.26.1 AS builder

# 国内加速环境
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off
ENV CGO_ENABLED=1

# 阿里云 Debian 源
RUN sed -i 's|deb.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's|security.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list.d/debian.sources

# 安装交叉编译器 + upx
RUN apt update && apt install -y --no-install-recommends \
    gcc-mingw-w64-x86-64 upx && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
# COPY . .
COPY ./cmd/pufo ./cmd/pufo
COPY ./ipkg/pufo ./ipkg/pufo

# 初始化依赖
RUN go mod init github.com/ymc-github/apkm || go mod tidy
RUN go mod tidy

#--------------------------
# 构建 Linux amd64
#--------------------------
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o pufo ./cmd/pufo
RUN upx -9 pufo 2>/dev/null || true

#--------------------------
# 构建 Windows amd64
#--------------------------
RUN CC=x86_64-w64-mingw32-gcc \
    GOOS=windows GOARCH=amd64 \
    go build -ldflags="-w -s" -o pufo.exe ./cmd/pufo
RUN upx -9 pufo.exe 2>/dev/null || true

#==============================================
# 导出双平台二进制
#==============================================
FROM scratch AS export
COPY --from=builder /app/pufo       ./pufo
COPY --from=builder /app/pufo.exe   ./pufo.exe

FROM debian:bookworm-slim AS runtime

# 阿里云源
RUN sed -i 's|deb.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list.d/debian.sources && \
    sed -i 's|security.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list.d/debian.sources

# 安装 PDF 转换依赖
RUN apt update && apt install -y --no-install-recommends \
    poppler-utils \
    fonts-wqy-zenhei \
    fontconfig && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
# COPY pufo .
COPY --from=builder /app/pufo  /usr/bin/pufo
RUN chmod +x /usr/bin/pufo
VOLUME ["/data"]

# ENTRYPOINT ["pufo"]
CMD ["pufo"]