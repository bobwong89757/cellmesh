package helpers

import "github.com/bobwong89757/gnbutils/yaml"

// MConfig 是全局的YAML配置工具实例
// 用于管理YAML格式的配置文件，提供配置的读取和缓存功能
var MConfig = &yaml.YamlUtil{
	KvsCache: make(map[string]map[string]string), // 配置缓存，键为配置文件名，值为键值对映射
}