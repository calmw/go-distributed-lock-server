## 镜像构建

``` shell
# build 
docker build -t distributed_lock_server:0.0.8 . 
# tag
docker tag distributed_lock_server:0.0.8 calmw/distributed_lock_server:0.0.8 
# push
docker push calmw/distributed_lock_server:0.0.8
```
