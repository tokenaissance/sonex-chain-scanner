package config

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainConfig struct {
	ChainId                 int64            `json:","`
	Name                    string           `json:","`
	Nodes                   map[string]int64 `json:","`
	MinConfirmations        int64            `json:",default=5"`
	MinDelayWaitMillisecond int64            `json:",default=2000"`
	MaxLogsNumber           int64            `json:",default=10000"`
	Multicall               *string          `json:",optional"`
	ETH                     *EthConfig       `json:",optional"`
}

type EthConfig struct {
	Name     string  `json:","`
	Symbol   string  `json:","`
	Decimals int64   `json:","`
	WETH     *string `json:",optional"`
}

func (c *ChainConfig) GetNode() *string {
	if len(c.Nodes) == 0 {
		return nil
	}
	totalWeight := linq.From(c.Nodes).SelectT(func(item linq.KeyValue) int64 { return item.Value.(int64) }).SumInts()
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randValue := rand.Int63n(totalWeight)
	var cumulativeWeight int64
	for key, weight := range c.Nodes {
		cumulativeWeight += weight
		if randValue < cumulativeWeight {
			return &key
		}
	}
	return nil
}

func (c *ChainConfig) GetClient() (*ethclient.Client, error) {
	node := c.GetNode()
	if node == nil {
		return nil, fmt.Errorf("%s", "Miss node configure")
	}
	client, err := ethclient.Dial(*node)
	if err != nil {
		return nil, err
	}
	return client, nil
}
