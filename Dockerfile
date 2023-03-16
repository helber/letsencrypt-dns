FROM golang:1.20
WORKDIR /go/src/github.com/helber/letsencrypt-dns
COPY . .
RUN go get ./...
RUN make && make install && make clean

FROM certbot/certbot:latest
# https://github.com/certbot/certbot
ENV GLIBC_VERSION=2.35-r0

## https://github.com/sgerrand/alpine-pkg-glibc
RUN apk --no-cache add ca-certificates wget bash && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$GLIBC_VERSION/glibc-$GLIBC_VERSION.apk && \
    apk add glibc-$GLIBC_VERSION.apk && rm -f glibc-$GLIBC_VERSION.apk

RUN apk  add --no-cache --virtual .certbot-deps \
    bash

# RUN wget -q https://mirror.openshift.com/pub/openshift-v4/clients/oc/4.4/linux/oc.tar.gz &&\
#     tar -zxf oc.tar.gz && \
#     mv oc /usr/local/bin/oc && \
#     chmod +x /usr/local/bin/oc && \
#     rm -Rf oc.tar.gz

COPY --from=0 /usr/local/bin/letsencrypt-dns /usr/local/bin/letsencrypt-dns
COPY --from=0 /usr/local/bin/letsencrypt-validate /usr/local/bin/letsencrypt-validate
COPY --from=0 /usr/local/bin/letsencrypt-cleanup /usr/local/bin/letsencrypt-cleanup
COPY --from=0 /usr/local/bin/checkcert  /usr/local/bin/checkcert
COPY --from=0 /usr/local/bin/oc-patch-route  /usr/local/bin/oc-patch-route
# Log directory inside container
RUN mkdir -p /var/log/letsencrypt-dns/

ENTRYPOINT [ "bash" ]
