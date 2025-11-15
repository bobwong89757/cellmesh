package memsd

import "time"

// Config 是memsd服务发现的配置结构
type Config struct {
	Address        string        // 服务发现服务器地址，格式为"host:port"
	RequestTimeout time.Duration // 请求超时时间
}

// DefaultConfig 返回默认的配置
// 默认地址为":8900"，超时时间为10秒
// 返回:
//   - *Config: 默认配置实例
func DefaultConfig() *Config {

	return &Config{
		Address:        ":8900",
		RequestTimeout: time.Second * 10,
	}
}
