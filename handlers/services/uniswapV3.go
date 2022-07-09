package services

import (
	"log"
	"math/big"
	"time"

	"github.com/Shuixingchen/year-service/utils"
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"

	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type UniswapV3Handler struct {
	Client *Client
}

type SwapParams struct {
	ChainID   uint64
	Token0    string
	Token1    string
	PoolFee   uint64
	SwapValue string
}

type Token struct {
	Address  string
	Dceimals uint8
	Name     string
	Symbol   string
}

func NewUniswapV3Handler() *UniswapV3Handler {
	return &UniswapV3Handler{Client: nil}
}

// ExactInput, 输入固定值token0, 得到token1
func (h *UniswapV3Handler) Swap(c *gin.Context) {
	p := &SwapParams{
		ChainID:   8001,
		Token0:    "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		Token1:    "0x0621d647cecbFb64b79E44302c1933cB4f27054d",
		PoolFee:   3000,
		SwapValue: "0.01",
	}
	client := NewClient(utils.Config.Nodes[int(p.ChainID)])
	h.Client = client

	wallet := InitWallet(utils.Config.PrivateKey)
	if wallet == nil {
		log.Fatal("init wallet failed")
	}
	token0 := h.getToken(uint(p.ChainID), p.Token0)
	token1 := h.getToken(uint(p.ChainID), p.Token1)

	// get pooladdress from FactoryContract and poolinfo from Uniswapv3Pool Contract
	pool, err := ConstructV3Pool(client.Clients[0], token0, token1, p.PoolFee)
	if err != nil {
		log.Fatal(err)
	}

	// 最大滑点0.01%
	slippageTolerance := coreEntities.NewPercent(big.NewInt(1), big.NewInt(1000))
	//after 5 minutes
	d := time.Now().Add(time.Minute * time.Duration(15)).Unix()
	deadline := big.NewInt(d)

	r, err := entities.NewRoute([]*entities.Pool{pool}, token0, token1)
	if err != nil {
		log.Fatal(err)
	}
	swapValue := FloatStringToBigInt(p.SwapValue, int(token0.Decimals()))
	trade, err := entities.FromRoute(r, coreEntities.FromRawAmount(token0, swapValue), coreEntities.ExactInput)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v %v\n", trade.Swaps[0].InputAmount.Quotient(), trade.Swaps[0].OutputAmount.Wrapped().Quotient())
	params, err := periphery.SwapCallParameters([]*entities.Trade{trade}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         wallet.PublicKey,
		Deadline:          deadline,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("calldata = 0x%x\n", params.Value.String())

	tx, err := TryTX(client.Clients[0], common.HexToAddress(ContractV3SwapRouterV1),
		swapValue, params.Calldata, wallet)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tx.Hash().String())
}

func (h *UniswapV3Handler) getToken(chainID uint, tokenAddr string) *coreEntities.Token {
	token := h.Client.GetToken(tokenAddr)
	return coreEntities.NewToken(chainID, common.HexToAddress(token.Address), uint(token.Dceimals), token.Name, token.Symbol)
}
