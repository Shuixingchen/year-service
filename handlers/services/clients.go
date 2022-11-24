package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Shuixingchen/year-service/contracts/artificial/erc20"
	"github.com/Shuixingchen/year-service/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	RPCClients []*rpc.Client
	Clients    []*ethclient.Client
}

func NewClient(nodes []utils.Node) *Client {
	clients := make([]*ethclient.Client, 0)
	rpcClients := make([]*rpc.Client, 0)
	for _, node := range nodes {
		c, err := rpc.DialContext(context.Background(), node.Addr)
		if err != nil {
			log.WithFields(log.Fields{"method:": "rpc.DialContext"}).Panic(err)
		}
		ec := ethclient.NewClient(c)
		clients = append(clients, ec)
		rpcClients = append(rpcClients, c)
	}
	return &Client{
		Clients:    clients,
		RPCClients: rpcClients,
	}
}

func (c *Client) GetToken(tokenAddr string) *Token {
	var token Token
	tc := c.GetTokenCaller(tokenAddr)
	token.Address = tokenAddr
	token.Name, _ = tc.Name(nil)
	token.Dceimals, _ = tc.Decimals(nil)
	token.Symbol, _ = tc.Symbol(nil)
	return &token
}

func (c *Client) GetTokenCaller(tokenAddr string) *erc20.Erc20 {
	contractAddress := common.HexToAddress(tokenAddr)
	tc, err := erc20.NewErc20(contractAddress, c.Clients[0])
	if err != nil {
		log.WithField("method", "NewErc20Caller").Error(err)
	}
	return tc
}

func (c *Client) GetLateBlock() (uint64, error) {
	return c.Clients[0].BlockNumber(context.Background())
}
func (c *Client) GetBlcokByNumber(number uint64) (*types.Block, error) {
	return c.Clients[0].BlockByNumber(context.Background(), big.NewInt(int64(number)))
}
func (c *Client) GetCode(contractAddr string) []byte {
	addr := common.HexToAddress(contractAddr)
	bytecode, err := c.Clients[0].CodeAt(context.Background(), addr, nil)
	if err != nil {
		fmt.Println(err)
	}
	return bytecode
}
