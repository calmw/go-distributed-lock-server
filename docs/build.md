## 镜像构建

``` shell
# build 
docker build -t distributed_lock_server:0.0.6 . 
# tag
docker tag distributed_lock_server:0.0.6 calmw/distributed_lock_server:0.0.6 
# push
docker push calmw/distributed_lock_server:0.0.6
```
