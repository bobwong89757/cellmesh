package kvconfig

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"reflect"
)

// doRaw 从服务发现系统获取配置值，如果不存在则使用默认值并自动写入
// 这是一个内部辅助函数，用于统一处理配置获取逻辑
// 参数:
//   - d: 服务发现实例，如果为nil则直接返回
//   - key: 配置项的键名
//   - defaultValue: 默认值，当配置不存在时使用
//   - ret: 指向目标变量的指针，用于接收配置值
func doRaw(d discovery.Discovery, key string, defaultValue, ret interface{}) {
	if d == nil {
		return
	}

	err := d.GetValue(key, ret)

	if err != nil && err.Error() == "value not exists" {

		reflect.Indirect(reflect.ValueOf(ret)).Set(reflect.ValueOf(defaultValue))
		// 默认值初始化
		d.SetValue(key, defaultValue)
	}

	return
}

// String 从服务发现系统获取字符串类型的配置值
// 如果配置不存在，会使用默认值并自动写入到服务发现系统
// 参数:
//   - d: 服务发现实例
//   - key: 配置项的键名
//   - defaultValue: 默认值
// 返回:
//   - ret: 配置值，如果不存在则返回默认值
func String(d discovery.Discovery, key string, defaultValue string) (ret string) {
	doRaw(d, key, defaultValue, &ret)
	return
}

// Int32 从服务发现系统获取int32类型的配置值
// 如果配置不存在，会使用默认值并自动写入到服务发现系统
// 参数:
//   - d: 服务发现实例
//   - key: 配置项的键名
//   - defaultValue: 默认值
// 返回:
//   - ret: 配置值，如果不存在则返回默认值
func Int32(d discovery.Discovery, key string, defaultValue int32) (ret int32) {
	doRaw(d, key, defaultValue, &ret)
	return
}

// Int64 从服务发现系统获取int64类型的配置值
// 如果配置不存在，会使用默认值并自动写入到服务发现系统
// 参数:
//   - d: 服务发现实例
//   - key: 配置项的键名
//   - defaultValue: 默认值
// 返回:
//   - ret: 配置值，如果不存在则返回默认值
func Int64(d discovery.Discovery, key string, defaultValue int64) (ret int64) {
	doRaw(d, key, defaultValue, &ret)
	return
}

// Bool 从服务发现系统获取bool类型的配置值
// 如果配置不存在，会使用默认值并自动写入到服务发现系统
// 参数:
//   - d: 服务发现实例
//   - key: 配置项的键名
//   - defaultValue: 默认值
// 返回:
//   - ret: 配置值，如果不存在则返回默认值
func Bool(d discovery.Discovery, key string, defaultValue bool) (ret bool) {
	doRaw(d, key, defaultValue, &ret)
	return
}
