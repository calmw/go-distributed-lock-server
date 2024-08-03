## 镜像构建

``` shell
# build 
docker build -t distributed_lock_server:0.1.0 . 
# tag
docker tag distributed_lock_server:0.1.0 calmw/distributed_lock_server:0.1.0 
# push
docker push calmw/distributed_lock_server:0.1.0
```
