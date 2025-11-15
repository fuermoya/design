package config

type System struct {
	WorkerId      uint16 `mapstructure:"worker-id" json:"worker-id" yaml:"worker-id"`
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`                               // 端口值
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
}
