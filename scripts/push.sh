#!/bin/sh

docker build -t kryptn/websub-listener:$1 .
docker push kryptn/websub-listener:$1
