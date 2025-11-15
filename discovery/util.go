package discovery

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// BytesToAny 将字节数组转换为指定的类型
// 支持的类型包括：int, float32, float64, bool, string，以及其他可通过JSON反序列化的类型
// 参数:
//   - data: 源字节数组
//   - dataPtr: 指向目标变量的指针，类型必须匹配
// 返回:
//   - error: 转换失败时返回错误信息
func BytesToAny(data []byte, dataPtr interface{}) error {

	switch ret := dataPtr.(type) {
	case *int:
		v, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*ret = int(v)
		return nil
	case *float32:
		v, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return err
		}
		*ret = float32(v)
		return nil
	case *float64:
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*ret = float64(v)
		return nil
	case *bool:
		v, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		*ret = v
		return nil
	case *string:
		*ret = string(data)
		return nil
	default:
		return json.Unmarshal(data, dataPtr)
	}
}

// ValueMetaToSlice 将ValueMeta数组转换为指定类型的切片
// 参数:
//   - pairs: ValueMeta数组，每个元素包含一个键值对
//   - dataPtr: 指向目标切片的指针，切片的元素类型必须与ValueMeta中的值类型匹配
// 返回:
//   - error: 转换失败时返回错误信息
func ValueMetaToSlice(pairs []ValueMeta, dataPtr interface{}) error {

	vdata := reflect.Indirect(reflect.ValueOf(dataPtr))

	elementCount := len(pairs)

	slice := reflect.MakeSlice(vdata.Type(), elementCount, elementCount)

	for i := 0; i < elementCount; i++ {

		sliceValue := reflect.New(slice.Type().Elem())

		err := BytesToAny(pairs[i].Value, sliceValue.Interface())

		if err != nil {
			return err
		}

		slice.Index(i).Set(sliceValue.Elem())
	}

	vdata.Set(slice)

	return nil

}

// AnyToBytes 将任意类型的数据转换为字节数组
// 支持的类型包括：int, int32, int64, uint32, uint64, float32, float64, bool, string
// 其他类型会通过JSON序列化
// 参数:
//   - data: 要转换的数据
//   - prettyPrint: 是否使用格式化输出（仅对JSON序列化有效）
// 返回:
//   - []byte: 转换后的字节数组
//   - error: 转换失败时返回错误信息
func AnyToBytes(data interface{}, prettyPrint bool) ([]byte, error) {

	switch v := data.(type) {
	case int, int32, int64, uint32, uint64, float32, float64, bool:
		return []byte(fmt.Sprint(data)), nil
	case string:
		return []byte(v), nil

	default:
		if prettyPrint {
			raw, err := json.MarshalIndent(data, "", "\t")
			if err != nil {
				return nil, err
			}

			return raw, nil
		} else {
			raw, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			return raw, nil
		}
	}
}
