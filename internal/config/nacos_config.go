package config

type NacosConfig struct {
	Host        string `json:","`
	Port        uint64 `json:","`
	NamespaceId string `json:","`
	DataId      string `json:","`
	Group       string `json:","`
}
