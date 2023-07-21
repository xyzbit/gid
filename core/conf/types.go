package conf

type Config struct {
	Server       *Server       `json:"server" yaml:"server"`
	DBConfig     *DBConfig     `json:"database" yaml:"database"`
	ConsulConfig *ConsulConfig `json:"consul" yaml:"consul"`
}

type DBConfig struct {
	Driver         string `json:"driver,omitempty"`
	Source         string `json:"source,omitempty"`
	MaxIdleConns   int    `json:"max_idle_conns,omitempty"`
	OpenConns      int    `json:"open_conns,omitempty"`
	MaxLifetimeSec int    `json:"max_lifetime_sec,omitempty"`
}

type ConsulConfig struct {
	Address string `json:"address,omitempty"`
	Scheme  string `json:"scheme,omitempty"`
}

type Server struct {
	Grpc *Grpc `json:"grpc,omitempty"`
}

type Grpc struct {
	Addr       string `json:"addr,omitempty"`
	Port       int    `json:"port,omitempty"`
	TimeoutSec int    `json:"timeout_sec,omitempty"`
}
