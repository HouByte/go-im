# go-im

# 项目结构
```shell
|   .gitignore #git 排除文件
|   go.mod     # 项目结构
|   main.go    # 服务器启动类
|   message.go # 消息处理 (策略模式)
|   server.go  # 服务器实现类
|   user.go    # 用户操作类
|   README.md  # 描述文件
```
# 服务器启动

```shell
go build im
./im
```
main.go 启动指定IP和端口
```go
package main

import "flag"

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器Port")
}
func main() {
	server := NewServer(serverIp, serverPort)
	server.Start()
}

```

# NC 客户端 
> 客户端就没有单独实现，利用Linux NC 客户端可以实现局域网聊天室
[nc](https://eternallybored.org/misc/netcat/)

```shell
nc [-hlnruz][-g<网关...>][-G<指向器数目>][-i<延迟秒数>][-o<输出文件>][-p<通信端口>][-s<来源位址>][-v...][-w<超时秒数>][主机名称][通信端口...]
```
连接服务器
```shell
 nc 127.0.0.1 8888
```
# 功能

|   命令   |     参数     |   介绍   | 示例        |
|:------:|:----------:|:------:|:----------|
|  who   |     无      | 在线用户列表 | who       |
|   bc   |    msg     |  广播消息  | bc hi!    |
| rename |    name    | 修改自己名称 | bc hi!    |
|   to   | receiveName msg |   私聊   | to tom hi |