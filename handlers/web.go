package handlers

import (
	"github.com/Shuixingchen/year-service/handlers/middleware"
	"github.com/Shuixingchen/year-service/handlers/services"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}
func (h *WebHandler) Handle() {
	r := gin.Default()
	r.Use(middleware.Cors())

	years := services.NewYearsHandler()
	uniswapV3 := services.NewUniswapV3Handler()
	wallet := services.NewWalletHandler()
	block := services.NewBlockHandler()
	worldcup := services.NewWorldCupHandler()
	evm := services.NewEVMHandler()

	versionRoute := r.Group("/v1")
	serviceRoute := versionRoute.Group("/years")
	serviceRoute.POST("/record", years.SaveRecord)

	v3 := r.Group("/v3")
	uniswapV3Route := v3.Group("/uniswap")
	uniswapV3Route.POST("/swap", uniswapV3.Swap)
	uniswapV3Route.POST("/quoter", uniswapV3.Quoter)
	uniswapV3Route.POST("/erc20/approve", uniswapV3.Approve)

	walletRoute := r.Group("/wallet")
	walletRoute.GET("/generate", wallet.Generate)

	blockRoute := r.Group("/block")
	blockRoute.GET("/latest", block.LatestBlock)

	worldcupRoute := r.Group("/worldcup")
	worldcupRoute.GET("/getallgame", worldcup.GetAllGames)

	evmRoute := r.Group("/evm")
	evmRoute.GET("/decompile", evm.Decompile)

	r.Run(":8080")
}
