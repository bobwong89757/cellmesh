package discovery

// ValueMeta 表示键值对配置的元数据信息
// 用于在服务发现系统中存储和传递配置数据
type ValueMeta struct {
	Key   string // 配置项的键名
	Value []byte // 配置项的值，以字节数组形式存储
}

// CheckerFunc 是健康检查函数的类型定义
// 返回两个字符串：output表示检查输出信息，status表示检查状态
type CheckerFunc func() (output, status string)

// Discovery 是服务发现接口，定义了服务注册、查询、配置管理等功能
// 实现该接口可以提供统一的服务发现能力，支持服务的自动发现和配置管理
type Discovery interface {

	// Register 注册一个服务到服务发现系统
	// 参数:
	//   - svc: 服务描述信息，包含服务名、ID、地址、端口等
	// 返回:
	//   - error: 注册失败时返回错误信息
	Register(*ServiceDesc) error

	// Deregister 从服务发现系统中注销指定的服务
	// 参数:
	//   - svcid: 服务的唯一标识ID
	// 返回:
	//   - error: 注销失败时返回错误信息
	Deregister(svcid string) error

	// Query 根据服务名称查询所有可用的服务实例
	// 参数:
	//   - name: 服务名称
	// 返回:
	//   - ret: 匹配的服务描述列表，如果没有找到则返回空列表
	Query(name string) (ret []*ServiceDesc)

	// RegisterNotify 注册服务变化通知通道
	// 当服务状态发生变化时，会通过返回的channel发送通知
	// 参数:
	//   - mode: 通知模式，如"add"表示服务添加通知
	// 返回:
	//   - ret: 用于接收通知的channel
	RegisterNotify(mode string) (ret chan struct{})

	// DeregisterNotify 解除服务变化通知的注册
	// 参数:
	//   - mode: 通知模式
	//   - c: 之前注册的通知channel
	DeregisterNotify(mode string, c chan struct{})

	// SetValue 在服务发现系统中设置一个配置值
	// 参数:
	//   - key: 配置项的键名
	//   - value: 配置项的值，可以是任意类型
	//   - optList: 可选的配置选项，如格式化选项等
	// 返回:
	//   - error: 设置失败时返回错误信息
	SetValue(key string, value interface{}, optList ...interface{}) error

	// GetValue 从服务发现系统中获取配置值并赋值到指定变量
	// 参数:
	//   - key: 配置项的键名
	//   - valuePtr: 指向目标变量的指针，用于接收配置值
	// 返回:
	//   - error: 获取失败时返回错误信息
	GetValue(key string, valuePtr interface{}) error

	// DeleteValue 从服务发现系统中删除指定的配置项
	// 参数:
	//   - key: 要删除的配置项的键名
	// 返回:
	//   - error: 删除失败时返回错误信息
	DeleteValue(key string) error
}

var (
	// Default 是默认的服务发现实例
	// 应用程序应该使用此实例进行服务注册、查询和配置管理
	Default Discovery
)
