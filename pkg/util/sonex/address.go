package sonex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ComputePoolAddress(factory common.Address, token0 common.Address, token1 common.Address, fee *big.Int) common.Address {
	if token0.Cmp(token1) > 0 {
		token0, token1 = token1, token0
	}

	poolData := append(common.LeftPadBytes(token0.Bytes(), 32), common.LeftPadBytes(token1.Bytes(), 32)...)
	poolData = append(poolData, common.LeftPadBytes(fee.Bytes(), 32)...)
	poolData = crypto.Keccak256(poolData)
	data := append([]byte{0xff}, factory.Bytes()...)
	data = append(data, poolData...)
	data = append(data, common.Hex2Bytes("e34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54")...)
	hash := crypto.Keccak256(data)

	return common.BytesToAddress(hash)
}
