package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	Nacos NacosConfig `json:","`
}

func (c *Config) GetSvcConfig(v any) error {
	// Load nacos configure
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig(c.Nacos.Host, c.Nacos.Port, constant.WithContextPath("/nacos")),
			},
			ClientConfig: constant.NewClientConfig(
				constant.WithNamespaceId(c.Nacos.NamespaceId),
				constant.WithTimeoutMs(5000),
				constant.WithNotLoadCacheAtStart(true),
			),
		},
	)
	if err != nil {
		return err
	}

	// Load service configure
	confString, err := client.GetConfig(vo.ConfigParam{
		DataId: c.Nacos.DataId,
		Group:  c.Nacos.Group,
	})
	if err != nil {
		return err
	}
	err = conf.LoadFromYamlBytes([]byte(confString), v)
	if err != nil {
		return err
	}
	return nil
}
