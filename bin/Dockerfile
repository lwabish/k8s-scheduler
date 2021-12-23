FROM alpine as helper

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache ca-certificates tzdata


FROM alpine as app

WORKDIR /app

COPY k8s-scheduler /app

COPY --from=helper /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=helper /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/k8s-scheduler"]