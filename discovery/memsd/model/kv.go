package model

import (
	"encoding/json"
	"github.com/bobwong89757/cellmesh/discovery"
	"io"
	"sort"
)

// ValueMeta 是KV存储中值的元数据
// 包含键、值、服务名（如果是服务描述）和认证令牌
type ValueMeta struct {
	Key     string // 键名
	Value   []byte // 值（字节数组）
	SvcName string // 服务名称，只有服务描述才有此字段
	Token   string // 认证令牌
}

// ErrDesc 是无效服务描述的占位符
var ErrDesc = discovery.ServiceDesc{Name: "invalid desc"}

// ValueAsServiceDesc 将值解析为服务描述
// 返回:
//   - *discovery.ServiceDesc: 解析后的服务描述，如果解析失败返回ErrDesc
func (self *ValueMeta) ValueAsServiceDesc() *discovery.ServiceDesc {

	var desc discovery.ServiceDesc
	err := json.Unmarshal(self.Value, &desc)
	if err != nil {
		return &ErrDesc
	}

	return &desc
}

var (
	valueByKey = map[string]*ValueMeta{}

	ValueDirty bool
)

func SetValue(key string, meta *ValueMeta) {
	ValueDirty = true
	valueByKey[key] = meta
}

func GetValue(key string) *ValueMeta {

	return valueByKey[key]
}

func DeleteValue(key string) *ValueMeta {
	ValueDirty = true
	ret := valueByKey[key]
	delete(valueByKey, key)

	return ret
}

func ValueCount() int {
	return len(valueByKey)
}

func VisitValue(callback func(*ValueMeta) bool) {
	for _, vmeta := range valueByKey {
		if !callback(vmeta) {
			return
		}
	}
}

// PersistFile 是持久化文件的结构
// 用于将内存中的数据保存到文件或从文件加载
type PersistFile struct {
	Version int          // 文件版本号
	Values  []*ValueMeta // 值列表
}

var (
	fileVersion = 1
)

func SaveValue(writer io.Writer) error {

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	var file PersistFile
	file.Version = fileVersion
	for _, vmeta := range valueByKey {
		file.Values = append(file.Values, vmeta)
	}

	sort.SliceStable(file.Values, func(i, j int) bool {

		return file.Values[i].Key < file.Values[j].Key
	})

	err := encoder.Encode(&file)

	if err != nil {
		return err
	}

	return nil
}

func LoadValue(reader io.Reader) error {

	decoder := json.NewDecoder(reader)

	var file PersistFile
	err := decoder.Decode(&file)
	if err != nil {
		return err
	}

	valueByKey = map[string]*ValueMeta{}

	for _, v := range file.Values {
		valueByKey[v.Key] = v
	}

	return nil
}
