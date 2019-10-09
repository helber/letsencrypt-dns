FROM golang:1.13
WORKDIR /go/src/github.com/helber/letsencrypt-dns
COPY . .
RUN make && make install && make clean


FROM certbot/certbot:latest
# https://github.com/certbot/certbot

## https://github.com/sgerrand/alpine-pkg-glibc
RUN apk --no-cache add ca-certificates wget bash && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk && \
    apk add glibc-2.28-r0.apk && rm -f glibc-2.28-r0.apk

RUN apk  add --no-cache --virtual .certbot-deps \
        bash
RUN wget https://github.com/openshift/origin/releases/download/v3.10.0/openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit.tar.gz
RUN tar -zxf openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit.tar.gz && \
    mv openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit/oc /usr/local/bin/oc && \
    chmod +x /usr/local/bin/oc && \
    rm -Rf openshift-origin-client-tools-v3.10.0-dd10d17-linux-64bit*

COPY --from=0 /usr/local/bin/letsencrypt-dns /usr/local/bin/letsencrypt-dns
COPY --from=0 /usr/local/bin/letsencrypt-validate /usr/local/bin/letsencrypt-validate
COPY --from=0 /usr/local/bin/letsencrypt-cleanup /usr/local/bin/letsencrypt-cleanup
COPY --from=0 /usr/local/bin/checkcert  /usr/local/bin/checkcert
COPY --from=0 /usr/local/bin/oc-patch-route  /usr/local/bin/oc-patch-route
# Log directory inside container
RUN mkdir -p /var/log/letsencrypt-dns/

ENTRYPOINT [ "bash" ]
