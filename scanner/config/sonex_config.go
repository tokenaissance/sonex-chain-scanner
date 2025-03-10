package config

type SonexConfig struct {
	Factory         string `json:","`
	PositionManager string `json:","`
	Quoter          string `json:",optional"`
	SwapRouter      string `json:",optional"`
}
