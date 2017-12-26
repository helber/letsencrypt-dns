FROM golang:1.9.2
WORKDIR /go/src/github.com/helber/letsencrypt-dns
COPY . .
RUN go get -d -v github.com/helber/letsencrypt-dns/...
RUN go get -d -v github.com/bobesa/go-domain-util/domainutil
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/letsencrypt-validate cmd/letsencrypt-validate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/letsencrypt-dns cmd/letsencrypt-dns/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/letsencrypt-cleanup cmd/letsencrypt-cleanup/main.go


FROM certbot/certbot:latest
# https://github.com/certbot/certbot

## https://github.com/sgerrand/alpine-pkg-glibc
RUN apk --no-cache add ca-certificates wget && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.26-r0/glibc-2.26-r0.apk && \
    apk add glibc-2.26-r0.apk && rm -f glibc-2.26-r0.apk

RUN apk  add --no-cache --virtual .certbot-deps \
        bash
COPY --from=0 /usr/local/bin/letsencrypt-dns /usr/local/bin/letsencrypt-dns
COPY --from=0 /usr/local/bin/letsencrypt-validate /usr/local/bin/letsencrypt-validate
COPY --from=0 /usr/local/bin/letsencrypt-cleanup /usr/local/bin/letsencrypt-cleanup

ENTRYPOINT [ "bash" ]
