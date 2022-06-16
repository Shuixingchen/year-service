package common

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySign(msg, signature, addr string) bool {

	// 1. 调用Ecrecover（椭圆曲线签名恢复）来检索签名者的公钥
	msgHash := crypto.Keccak256Hash([]byte(msg))
	sig, _ := hex.DecodeString(signature)
	sigPublicKey, err := crypto.SigToPub(msgHash.Bytes(), sig)
	if err != nil {
		log.Fatal(err)
	}
	// 2.公钥得到地址
	signAddr := crypto.PubkeyToAddress(*sigPublicKey).Hex()

	// 3.对比两个地址是否一致
	matches := strings.EqualFold(addr, signAddr)
	return matches
}
