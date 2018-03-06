FROM golang:1.10
WORKDIR /go/src/github.com/helber/letsencrypt-dns
COPY . .
RUN go get -d -v github.com/helber/letsencrypt-dns/...
RUN go get -d -v github.com/bobesa/go-domain-util/domainutil
RUN make && make install && make clean


FROM certbot/certbot:latest
# https://github.com/certbot/certbot

## https://github.com/sgerrand/alpine-pkg-glibc
RUN apk --no-cache add ca-certificates wget && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.26-r0/glibc-2.26-r0.apk && \
    apk add glibc-2.26-r0.apk && rm -f glibc-2.26-r0.apk

RUN apk  add --no-cache --virtual .certbot-deps \
        bash
# RUN wget https://github.com/openshift/origin/releases/download/v3.7.1/openshift-origin-client-tools-v3.7.1-ab0f056-linux-64bit.tar.gz
# RUN tar -zxf openshift-origin-client-tools-v3.7.1-ab0f056-linux-64bit.tar.gz && \
#     mv openshift-origin-client-tools-v3.7.1-ab0f056-linux-64bit/oc /usr/local/bin/oc && \
#     chmod +x /usr/local/bin/oc && \
#     rm -Rf openshift-origin-client-tools-v3.7.1-ab0f056-linux-64bit*

COPY --from=0 /usr/local/bin/letsencrypt-dns /usr/local/bin/letsencrypt-dns
COPY --from=0 /usr/local/bin/letsencrypt-validate /usr/local/bin/letsencrypt-validate
COPY --from=0 /usr/local/bin/letsencrypt-cleanup /usr/local/bin/letsencrypt-cleanup
COPY --from=0 /usr/local/bin/checkcert  /usr/local/bin/checkcert
# Log directory inside container
RUN mkdir -p /var/log/letsencrypt-dns/

ENTRYPOINT [ "bash" ]
