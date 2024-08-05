## 安装protoc以及插件

- 1 安装protoc命令 ``` brew install protobuf ```

- 2 安装protobuf插件 ``` go install google.golang.org/protobuf/cmd/protoc-gen-go@latest ```
- 3 安装protobuf插件 ``` go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest ```,
  生成gRPC相关代码需要安装grpc-go相关的插件protoc-gen-go-grpc

## 使用protoc命令生成go代码

- 目录中的子目录需要事先创建

``` shell
# 将lock.proto放到service目录，并在该目录下执行下面命令
protoc --go_out=../service --go_opt=paths=source_relative \
    --go-grpc_out=../service --go-grpc_opt=paths=source_relative \
    --proto_path ../service \
    lock.proto
```