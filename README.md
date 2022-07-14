# year-service
区块链相关api
包括btc 上链数据
uniswapV3 接口

```shell
go run main.go

```
### 接口列表
| 接口 | 请求方式 | 介绍 |
| ------ | ------ | ------ |
| /v1/years/record | POST | metaMask签名后的数据 |
| /v3/uniswap/swap | POST | UniswapV3兑换接口，ExcuteInput |