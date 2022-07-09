package services

import (
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	sdkutils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  common.Address
}

func (w Wallet) PubkeyStr() string {
	return w.PublicKey.String()
}

func InitWallet(privateHexKeys string) *Wallet {
	if privateHexKeys == "" {
		return nil
	}
	privateKey, err := crypto.HexToECDSA(privateHexKeys)
	if err != nil {
		return nil
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

func GetPoolAddress(client *ethclient.Client, token0, token1 common.Address, fee *big.Int) (common.Address, error) {
	f, err := contract.NewUniswapv3Factory(common.HexToAddress(ContractV3Factory), client)
	if err != nil {
		return common.Address{}, err
	}
	poolAddr, err := f.GetPool(nil, token0, token1, fee)
	if err != nil {
		return common.Address{}, err
	}
	if poolAddr == (common.Address{}) {
		return common.Address{}, errors.New("pool is not exist")
	}

	return poolAddr, nil
}

func ConstructV3Pool(client *ethclient.Client, token0, token1 *coreEntities.Token, poolFee uint64) (*entities.Pool, error) {
	poolAddress, err := GetPoolAddress(client, token0.Address, token1.Address, new(big.Int).SetUint64(poolFee))
	if err != nil {
		return nil, err
	}

	contractPool, err := contract.NewUniswapv3Pool(poolAddress, client)
	if err != nil {
		return nil, err
	}

	liquidity, err := contractPool.Liquidity(nil)
	if err != nil {
		return nil, err
	}

	slot0, err := contractPool.Slot0(nil)
	if err != nil {
		return nil, err
	}

	pooltick, err := contractPool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return nil, err
	}

	feeAmount := constants.FeeAmount(poolFee)
	ticks := []entities.Tick{
		{
			Index: entities.NearestUsableTick(sdkutils.MinTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
		{
			Index: entities.NearestUsableTick(sdkutils.MaxTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
	}

	// create tick data provider
	p, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	if err != nil {
		return nil, err
	}

	return entities.NewPool(token0, token1, constants.FeeAmount(poolFee),
		slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}

func GetTokenInstance(client *ethclient.Client, tokenAddr string) {

	WMATIC = coreEntities.NewToken(PolygonChainID, common.HexToAddress(WMaticAddr), 18, "Matic", "Matic Network(PolyGon)")
}

func FloatStringToBigInt(amount string, decimals int) *big.Int {
	fAmount, _ := new(big.Float).SetString(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(decimals))).Int(nil)
	return fi
}
