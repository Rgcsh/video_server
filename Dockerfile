FROM gocv/opencv:4.5.4 AS build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://mirrors.aliyun.com/goproxy/,direct
WORKDIR  /release

ADD . .
RUN go mod tidy && go mod vendor
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o video_server main.go

FROM gocv/opencv:4.5.4

ENV LANG C.UTF-8

WORKDIR /data

COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build /release/video_server .

# 设置时区
RUN mkdir log video_file
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone \
    && cp /etc/apt/sources.list /etc/apt/sources.list.bak \
    && sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y vim

EXPOSE 7069
CMD ["./video_server"]
