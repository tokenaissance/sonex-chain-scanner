package config

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
