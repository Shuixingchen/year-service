package services

import (
	"math/big"
	"net/http"
	"time"

	"github.com/Shuixingchen/year-service/utils"
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"

	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UniswapV3Handler struct {
	Client     *Client
	PrivateKey string
}

type SwapParams struct {
	ChainID     uint64
	Token0      string
	Token1      string
	PoolFee     uint64
	SwapValue   string
	MaxSlippage uint64
}
type ApproveParams struct {
	Token   string
	Account string
	Amount  *big.Int
}

type Token struct {
	Address  string
	Dceimals uint8
	Name     string
	Symbol   string
}

const (
	SportToken     = "0xE50498ec2aFA8121c763A7f06AA55b19f86Af013"
	GoveranceToken = "0x7Fe7Bc8a6bd64dFA54422ac96B39b2AB13e5d9cD"
	WMATIC         = "0x9c3C9283D3e44854697Cd22D3Faa240Cfb032889"
)
const (
	ContractV3Factory            = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	ContractV3SwapRouterV1       = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	ContractV3SwapRouterV2       = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"
	ContractV3NFTPositionManager = "0xC36442b4a4522E871399CD717aBDD847Ab11FE88"
	ContractV3Quoter             = "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
)

func NewUniswapV3Handler() *UniswapV3Handler {
	return &UniswapV3Handler{Client: nil, PrivateKey: utils.Config.PrivateKey}
}

// ExactInput, 输入固定值token0, 得到token1
func (h *UniswapV3Handler) Swap(c *gin.Context) {
	// p := &SwapParams{
	// 	ChainID:     80001,
	// 	Token0:      WMATIC,
	// 	Token1:      SportToken,
	// 	PoolFee:     3000, // 最终除以1000000, 0.3%
	// 	SwapValue:   "0.01",
	// 	MaxSlippage: 10, // 最终也是除以1000000
	// }
	var p SwapParams
	c.BindJSON(&p)

	if len(utils.Config.Nodes[int(p.ChainID)]) == 0 {
		log.Panic("need config chainid: ", p.ChainID)
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
		log.WithField("method", "ConstructV3Pool").Fatal(err)
	}

	slippageTolerance := coreEntities.NewPercent(big.NewInt(int64(p.MaxSlippage)), big.NewInt(1000000))
	//after 5 minutes
	d := time.Now().Add(time.Minute * time.Duration(15)).Unix()
	deadline := big.NewInt(d)

	r, err := entities.NewRoute([]*entities.Pool{pool}, token0, token1)
	if err != nil {
		log.WithField("method", "NewRoute").Fatal(err)
	}
	swapValue := FloatStringToBigInt(p.SwapValue, int(token0.Decimals()))
	trade, err := entities.FromRoute(r, coreEntities.FromRawAmount(token0, swapValue), coreEntities.ExactInput)
	if err != nil {
		log.WithField("method", "FromRoute").Fatal(err)
	}
	log.Printf("%v %v\n", trade.Swaps[0].InputAmount.Quotient(), trade.Swaps[0].OutputAmount.Wrapped().Quotient())
	params, err := periphery.SwapCallParameters([]*entities.Trade{trade}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         wallet.PublicKey,
		Deadline:          deadline,
	})
	if err != nil {
		log.WithField("method", "SwapCallParameters").Fatal(err)
	}

	tx, err := SendTX(client.Clients[0], common.HexToAddress(ContractV3SwapRouterV1),
		swapValue, params.Calldata, wallet)
	if err != nil {
		log.WithField("method", "SendTX").Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"txHash": tx.Hash().Hex(),
	})
}

func (h *UniswapV3Handler) getToken(chainID uint, tokenAddr string) *coreEntities.Token {
	token := h.Client.GetToken(tokenAddr)
	return coreEntities.NewToken(chainID, common.HexToAddress(token.Address), uint(token.Dceimals), token.Name, token.Symbol)
}

// erc20 第一次swap，需要Approve给
func (h *UniswapV3Handler) Approve(c *gin.Context) {
	var p ApproveParams
	c.BindJSON(&p)
	tc := h.Client.GetTokenCaller(p.Token)
	opts := h.GetBindOption()
	tx, err := tc.Approve(opts, common.HexToAddress(p.Account), p.Amount)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"txHash": tx.Hash().Hex(),
	})
}

func (h *UniswapV3Handler) GetBindOption() *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(h.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	return auth
}
