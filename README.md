# websub-listener

This service builds webhooks and manages subscriptions to a [WebSub Hub](https://pubsubhubbub.appspot.com/) and sends a slack message through an incoming webhook



## usage

Create a config file, and run with docker

```
docker run -p 8080:8080 -v $PWD/config.toml:/etc/websub/config.toml kryptn/websub-listener:v0.0.9
```