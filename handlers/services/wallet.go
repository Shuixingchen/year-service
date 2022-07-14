package services

import (
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

func (h *WalletHandler) Generate(c *gin.Context) {

}
