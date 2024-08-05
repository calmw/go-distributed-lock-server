## 镜像构建

``` shell
# build 
docker build -t distributed_lock_server:0.1.2 . 
# tag
docker tag distributed_lock_server:0.1.2 calmw/distributed_lock_server:0.1.2 
# push
docker push calmw/distributed_lock_server:0.1.2
```
