package core

import "time"

var (
	WatchReStartInterval time.Duration = 3 * time.Second // 监听重启间隔; 如果panic, 重启间隔
	WatchRetryInterval   time.Duration = 3 * time.Second // 监听重试次数; 如果更新状态失败，重试间隔
)

type Status struct {
	ModuleID   int32 // 模块id
	ModuleName string
	LastID     int64 // 最后生成的id
}

type Config struct {
	// ModuleID             int32 // 模块id
	ModuleName    string `json:"module_name,omitempty"`
	Desc          string `json:"desc,omitempty"`
	IDKind        int8   `json:"id_kind,omitempty"`        // id类型
	LoadNum       int32  `json:"load_num,omitempty"`       // 每批生成数量
	LoadThreshold int32  `json:"load_threshold,omitempty"` // 加载阀值

	IncrFirstID int64 `json:"incr_first_id,omitempty"` // 第一个id
}

type StatusLoader interface {
	// 获取所有生成器生成状态
	AllStatus() ([]*Status, error)
	// 获取生成器生成状态
	Status(name string) (*Status, error)
	// 更新生成器生成状态
	SaveStatus(name string, s *Status) error
}

type ConfigLoader interface {
	// 获取所有生成器生成配置
	AllConfigs() ([]*Config, error)
	// 获取生成器生成配置
	Config(name string) (*Config, error)
	// 更新生成器生成配置
	SaveConfig(name string, s *Config) error
}
