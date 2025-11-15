package discovery

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ServiceDesc 表示注册到服务发现系统的服务描述信息
// 包含了服务的所有元数据，用于服务发现、路由和负载均衡
type ServiceDesc struct {
	Name string            // 服务名称，用于标识服务类型，如"game"、"login"等
	ID   string            // 服务的唯一标识ID，格式通常为"服务名#索引@分组"，在所有服务中必须唯一
	Host string            // 服务所在的主机地址，可以是IP地址或域名
	Port int               // 服务监听的端口号
	Tags []string          // 服务的分类标签，用于服务筛选和分组
	Meta map[string]string // 服务的元数据配置，存储额外的服务信息，如分组、索引、外网地址等
}

// Equals 比较两个服务描述是否完全相等
// 比较所有字段，包括ID、端口、名称、主机、标签和元数据
// 参数:
//   - sd: 要比较的另一个服务描述
//
// 返回:
//   - bool: 如果所有字段都相等返回true，否则返回false
func (self *ServiceDesc) Equals(sd *ServiceDesc) bool {

	if sd.ID != self.ID {
		return false
	}

	if sd.Port != self.Port {
		return false
	}

	if sd.Name != self.Name {
		return false
	}

	if sd.Host != self.Host {
		return false
	}

	if !reflect.DeepEqual(self.Tags, sd.Tags) {
		return false
	}

	if !reflect.DeepEqual(self.Meta, sd.Meta) {
		return false
	}

	return true
}

// ContainTags 检查服务描述是否包含指定的标签
// 参数:
//   - tag: 要检查的标签名称
//
// 返回:
//   - bool: 如果包含该标签返回true，否则返回false
func (self *ServiceDesc) ContainTags(tag string) bool {
	for _, libtag := range self.Tags {
		if libtag == tag {
			return true
		}
	}

	return false
}

// SetMeta 设置服务的元数据项
// 如果Meta字段为nil，会自动初始化
// 参数:
//   - key: 元数据的键名
//   - value: 元数据的值
func (self *ServiceDesc) SetMeta(key, value string) {
	if self.Meta == nil {
		self.Meta = make(map[string]string)
	}

	self.Meta[key] = value
}

// GetMeta 获取服务的元数据值
// 参数:
//   - name: 元数据的键名
//
// 返回:
//   - string: 元数据的值，如果不存在则返回空字符串
func (self *ServiceDesc) GetMeta(name string) string {
	if self.Meta == nil {
		return ""
	}

	return self.Meta[name]
}

// GetMetaAsInt 获取服务的元数据值并转换为整数
// 参数:
//   - name: 元数据的键名
//
// 返回:
//   - int: 转换后的整数值，如果转换失败或不存在则返回0
func (self *ServiceDesc) GetMetaAsInt(name string) int {
	v, err := strconv.ParseInt(self.GetMeta(name), 10, 64)
	if err != nil {
		return 0
	}

	return int(v)
}

// Address 返回服务的完整地址，格式为"host:port"
// 返回:
//   - string: 服务的网络地址
func (self *ServiceDesc) Address() string {
	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

// String 返回服务描述的字符串表示
// 包含服务ID、主机、端口和元数据信息
// 返回:
//   - string: 格式化的服务描述字符串
func (self *ServiceDesc) String() string {
	var sb strings.Builder
	if len(self.Meta) > 0 {

		sb.WriteString("meta: [ ")
		for key, value := range self.Meta {
			sb.WriteString(key)
			sb.WriteString("=")
			sb.WriteString(value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%s host: %s port: %d %s", self.ID, self.Host, self.Port, sb.String())
}

// FormatString 返回格式化的服务描述字符串
// 与String()方法类似，但使用固定宽度格式，便于对齐显示
// 元数据按键名排序输出
// 返回:
//   - string: 格式化的服务描述字符串
func (self *ServiceDesc) FormatString() string {

	var sb strings.Builder
	if len(self.Meta) > 0 {

		type pair struct {
			key   string
			value string
		}

		var pairs []pair

		for key, value := range self.Meta {
			pairs = append(pairs, pair{key, value})
		}

		sort.Slice(pairs, func(i, j int) bool {

			return pairs[i].key < pairs[j].key
		})

		sb.WriteString("meta: [ ")
		for _, kv := range pairs {
			sb.WriteString(kv.key)
			sb.WriteString("=")
			sb.WriteString(kv.value)
			sb.WriteString(" ")
		}
		sb.WriteString("]")
	}

	return fmt.Sprintf("%25s host: %15s port: %5d %s", self.ID, self.Host, self.Port, sb.String())
}
