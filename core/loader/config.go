package loader

import (
	"encoding/json"
	"fmt"

	"github.com/xyzbit/gid/core"
	"github.com/xyzbit/gid/core/conf"

	consulAPI "github.com/hashicorp/consul/api"
)

type ConfigLoader struct {
	cli *consulAPI.Client
}

func NewConfigLoader(cc *conf.ConsulConfig) core.ConfigLoader {
	return &ConfigLoader{
		cli: NewConsulCli(cc),
	}
}

func NewConsulCli(cc *conf.ConsulConfig) *consulAPI.Client {
	c := consulAPI.DefaultConfig()
	c.Address = cc.Address
	c.Scheme = cc.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	return cli
	// r := consul.New(cli, consul.WithHealthCheck(false))
}

// 实现 ConfigLoader 接口
func (l *ConfigLoader) AllConfigs() ([]*core.Config, error) {
	kvs, metadata, err := l.cli.KV().List("modules", nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("AllConfigs metadata %+v", metadata)

	rets := make([]*core.Config, 0, len(kvs))
	for _, kv := range kvs {
		ret := &core.Config{}
		if err := json.Unmarshal(kv.Value, ret); err != nil {
			fmt.Printf("Error AllConfigs %+v", err)
			continue
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

func (l *ConfigLoader) Config(name string) (*core.Config, error) {
	kv, metadata, err := l.cli.KV().Get(fmt.Sprintf("modules/%s", name), nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Config metadata %+v", metadata)

	ret := &core.Config{}
	if err := json.Unmarshal(kv.Value, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (l *ConfigLoader) SaveConfig(name string, s *core.Config) error {
	cfg, err := json.Marshal(s)
	if err != nil {
		return err
	}
	metadata, err := l.cli.KV().Put(&consulAPI.KVPair{
		Key:   fmt.Sprintf("modules/%s", name),
		Value: cfg,
	}, nil)
	if err != nil {
		return err
	}
	fmt.Printf("SaveConfig metadata %+v", metadata)
	return nil
}
