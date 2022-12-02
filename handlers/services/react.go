package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Shuixingchen/year-service/utils"
	"github.com/gin-gonic/gin"
)

type React struct {
}
type LoginParam struct {
	Mobile string `json:"mobile"`
	Code   string `json:"code"`
}

const LoginToken = "abcdefghi"

func NewReactHandler() *React {
	return &React{}
}

func (h *React) DoLogin(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	var param LoginParam
	json.Unmarshal(data, &param)
	if len(param.Mobile) > 0 && param.Code == "123456" {
		utils.Response(c, 200, LoginToken)
		return
	}
	utils.Response(c, 403, "wrong code")
}
