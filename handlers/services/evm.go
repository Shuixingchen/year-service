package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shuixingchen/year-service/utils"
	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type EVM struct {
	Client *Client
}

type Instruction struct {
	PC  string // 指令的地址
	Op  string // opcode
	Arg string // 参数
}

func NewEVMHandler() *EVM {
	return &EVM{}
}

func (h *EVM) Decompile(c *gin.Context) {
	chainId := c.Query("chain_id")
	constractAddr := c.Query("contract_addr")
	id, _ := strconv.Atoi(chainId)
	if len(utils.Config.Nodes[int(id)]) == 0 {
		log.Panic("need config chainid: ", id)
	}
	client := NewClient(utils.Config.Nodes[int(id)])
	byteCode := client.GetCode(constractAddr)
	ins := decompile(byteCode)
	res := ""
	for _, v := range ins {
		res = res + fmt.Sprintf("%s  %s  %s \n", v.PC, v.Op, v.Arg)
	}
	utils.Response(c, http.StatusOK, res)
}

func decompile(byteCode []byte) []*Instruction {
	instructions := make([]*Instruction, 0)
	it := asm.NewInstructionIterator(byteCode)
	for it.Next() {
		var inst Instruction
		inst.PC = fmt.Sprintf("%05x", it.PC())
		inst.Op = fmt.Sprintf("%v", it.Op())
		if it.Arg() != nil && 0 < len(it.Arg()) {
			inst.Arg = fmt.Sprintf("%x", it.Arg())
		}
		instructions = append(instructions, &inst)
	}
	return instructions
}
