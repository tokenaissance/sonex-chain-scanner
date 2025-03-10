package config

type ChainConfig struct {
	ChainId                  int64             `json:","`
	Name                     string            `json:","`
	Nodes                    map[string]int64  `json:","`
	MinConfirmations         int64             `json:",default=5"`
	MinDelayWaitMilliseconds int64             `json:",default=2000"`
	MaxDelayWaitMilliseconds int64             `json:",default=10000"`
	MaxBlocksNumber          int64             `json:",default=1000"`
	Multicall                *string           `json:",optional"`
	NativeCoin               *NativeCoinConfig `json:",optional"`
}

type NativeCoinConfig struct {
	Name     string `json:","`
	Symbol   string `json:","`
	Decimals int64  `json:","`
	WrapCoin string `json:",optional"`
}
