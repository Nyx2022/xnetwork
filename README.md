# xnetwork

## Intro

tcp/udp client for scanner

应用于扫描器场景下的tcp/udp基础库。

1. client

   - 精准的http client配置：目前支持支持9项
   - 失败重试
   - limiter：qps限制
   - 超时
   - tls
2. request

   - GetRaw：请求内容
3. response

   - SourceAddr：源地址信息
   - DestinationAddr：目的地址信息
   - GetLatency：从发起请求到捕获响应的持续时间
   - GetRaw：响应内容
4. responseMiddleware：响应获取后，对响应的处理
   - 目前只有 debug 模式下需要的打印功能
5. debug模式：debug模式下将打印请求和响应完整信息
6. 完整的 tcp test server

## Install

```
go get github.com/xiecat/xnetwork
```

## Demo

```go
client := NewClient()
ctx := context.Background()

addr := &input.ServiceAsset{
		Host:    "127.0.0.1",
		Port:    "3306",
		Network: "tcp",
	}
// 目标连接
err := client.Dial(ctx, addr)
if err != nil {
  return
}
// 构造请求
req := &Request{
		Raw: []byte("FIRST"),
}
// 发起请求
resp, err := client.Do(req)
```

## Todo

- errorHook

## Ref

- https://github.com/xiecat/xhttp
