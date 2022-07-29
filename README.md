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


### abigen编译合约
```shell
abigen --sol contract/QuoterV2.sol --pkg uniswap --out artificial/uniswap/quoterV2.go
abigen --sol contracts/erc721.sol --pkg erc721 --out artificial/erc721/erc721.go
```