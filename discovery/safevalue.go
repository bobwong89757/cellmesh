package discovery

import (
	"fmt"
	"github.com/bobwong89757/cellnet/util"
	"reflect"
)

// PackedValueSize 定义KV存储中单个分片的最大大小
// 注意：由于底层使用JSON转base64编码，base64编码后的数据比原二进制大约33%，
// 所以实际二进制数据不到512K就会达到限制，这里设置为300KB以确保安全
const (
	PackedValueSize = 300 * 1024 // 单个分片的最大大小，单位：字节
)

// rawGetter 是获取原始值的接口
// 用于支持大值分片存储和读取的内部接口
type rawGetter interface {
	// GetRawValue 获取指定键的原始字节值
	// 参数:
	//   - key: 配置项的键名
	// 返回:
	//   - []byte: 原始字节数据
	//   - error: 获取失败时返回错误信息
	GetRawValue(key string) ([]byte, error)
	
	// GetValueDirect 直接获取配置值并赋值到指定变量
	// 参数:
	//   - key: 配置项的键名
	//   - valuePtr: 指向目标变量的指针
	// 返回:
	//   - error: 获取失败时返回错误信息
	GetValueDirect(key string, valuePtr interface{}) error
}

// getMultiKey 获取大值分片存储的所有键名列表
// 大值会被分割成多个分片，使用key, key.1, key.2...的格式存储
// 参数:
//   - sd: 实现了rawGetter接口的对象
//   - key: 主键名
// 返回:
//   - ret: 所有分片键名的列表，包括主键和分片键
func getMultiKey(sd rawGetter, key string) (ret []string) {

	mainKey := key

	ret = append(ret, mainKey)

	for i := 1; ; i++ {

		key = fmt.Sprintf("%s.%d", mainKey, i)

		_, err := sd.GetRawValue(key)
		if err != nil && err.Error() == "value not exists" {
			return
		}

		ret = append(ret, key)
	}

}

// SafeSetValue 安全地设置配置值，支持大值分片存储和压缩
// 当值较大时，会自动分割成多个分片存储（key, key.1, key.2...）
// 参数:
//   - sd: 服务发现实例
//   - key: 配置项的键名
//   - value: 配置项的值，compress为true时必须是[]byte类型
//   - compress: 是否启用压缩，启用后会对数据进行压缩后再存储
// 返回:
//   - error: 设置失败时返回错误信息
func SafeSetValue(sd Discovery, key string, value interface{}, compress bool) error {
	if compress {
		cData, err := util.CompressBytes(value.([]byte))
		if err != nil {
			return err
		}

		if len(cData) >= PackedValueSize {

			for _, multiKey := range getMultiKey(sd.(rawGetter), key) {

				err := sd.DeleteValue(multiKey)
				if err != nil {
					fmt.Printf("delete kv error, %s\n", err)
				}
			}

			var pos = PackedValueSize

			err = sd.SetValue(key, cData[:pos])
			if err != nil {
				return err
			}

			index := 1
			for len(cData)-pos > PackedValueSize {

				multiKey := fmt.Sprintf("%s.%d", key, index)
				err = sd.SetValue(multiKey, cData[pos:pos+PackedValueSize])
				if err != nil {
					return err
				}
				pos += PackedValueSize
				index++
			}

			if len(cData)-pos > 0 {
				multiKey := fmt.Sprintf("%s.%d", key, index)
				err = sd.SetValue(multiKey, cData[pos:])
				if err != nil {
					return err
				}
			}

			return nil

		} else {
			return sd.SetValue(key, cData)
		}

	} else {
		return sd.SetValue(key, value)
	}
}

// SafeGetValue 安全地获取配置值，支持大值分片读取和解压缩
// 如果值被分片存储，会自动合并所有分片；如果被压缩，会自动解压
// 参数:
//   - sd: 服务发现实例，必须实现rawGetter接口
//   - key: 配置项的键名
//   - valuePtr: 指向目标变量的指针，用于接收配置值
//   - decompress: 是否启用解压缩，如果存储时使用了压缩，这里必须为true
// 返回:
//   - error: 获取失败时返回错误信息
func SafeGetValue(sd Discovery, key string, valuePtr interface{}, decompress bool) error {

	rg := sd.(rawGetter)

	var (
		finalData []byte
		err       error
	)

	if decompress {

		var data []byte
		for _, multiKey := range getMultiKey(rg, key) {

			var partData []byte
			err := rg.GetValueDirect(multiKey, &partData)
			if err != nil {
				return err
			}

			data = append(data, partData...)
		}

		finalData, err = util.DecompressBytes(data)

		if err != nil {
			return err
		}

		reflect.ValueOf(valuePtr).Elem().Set(reflect.ValueOf(finalData))
	} else {
		err = rg.GetValueDirect(key, &finalData)

		if err != nil {
			return err
		}
	}

	return nil
}
