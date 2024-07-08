# Mutex互斥锁在银行交易场景下的使用

> 本项目是一个简单的RPC服务端和客户端的例子，通过互斥锁保护共享资源，模拟银行账户的存取款操作。

## 1. 项目背景

在银行交易场景下，多个线程对同一个账户进行存取款操作时，可能会出现数据不一致的情况。为了保证数据的一致性，需要使用互斥锁来保护共享资源。

客户端的日志是乱序的，建议查看服务端的日志。

### server

- 启动RPC服务端
- 通过互斥锁保护模拟账户的存取款操作【这里是用的全局变量模拟】
- 禁止交易额余额为负数

### client

- 启动RPC客户端
- 模拟多个客户端对同一个账户进行存取款操作【gorouting 并发模拟账户在多地发生交易】
- 通过RPC调用server端的存取款操作

## 2. 本地运行

```shell
cd SafeGoroutine
go mod tidy

# server
go run server/main.go

# client
go run client/main.go 
```
