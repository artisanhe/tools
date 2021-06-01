# README

## 依赖安装

```
git clone https://github.com/apache/rocketmq-client-cpp
mkdir -p /usr/local/include/rocketmq/
cp rocketmq-client-cpp/include/* /usr/local/include/rocketmq
sh build.sh
cp bin/librocketmq.dylib /usr/local/lib
go get github.com/apache/rocketmq-client-go
```