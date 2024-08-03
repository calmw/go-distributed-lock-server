# go-distributed-lock

- golang 实现的基于GRPC的分布式锁

## 安装

- 服务端，可以直接使用docker部署:calmw/distributed_lock_serve:latest

```shell
distributed_lock:
    image: calmw/distributed_lock_server:0.0.9
    container_name: distributed_lock_server
    restart: always
    environment:
      - LISTEN_ADDR=0.0.0.0:6000
```

- 客户端,参考 [example](example)
