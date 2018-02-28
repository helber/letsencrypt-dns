#!/bin/bash
PW=`pwd`
DIR=`dirname $PW`
echo $DIR
docker run \
    -it \
    --entrypoint=/bin/bash \
    -e LINODE_API_KEY=$LINODE_API_KEY \
    -e CF_API_EMAIL=$CF_API_EMAIL \
    -e CF_API_KEY=$CF_API_KEY \
    -v $DIR:/app helber/letsencrypt-dns

