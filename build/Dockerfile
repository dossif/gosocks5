FROM alpine:3.18.3
RUN apk add --no-cache ca-certificates
WORKDIR /home
ADD ./gosocks5 .
ENTRYPOINT [ "sh", "-c", "./gosocks5"]
