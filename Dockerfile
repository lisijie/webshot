FROM golang:1.14-alpine as builder
WORKDIR /root
COPY . /root/
RUN go build -mod=vendor -o webshot

FROM chromedp/headless-shell:latest as prod
# 安装中文字体和supervisor
RUN mkdir /app && apt-get update \
    && apt-get install -y --no-install-recommends \
        fonts-wqy-zenhei \
        supervisor \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=0 /root/webshot /app/
COPY ./supervisord.conf /etc/supervisor/

WORKDIR /app
EXPOSE 80

ENTRYPOINT ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]