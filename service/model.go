package service

var (
	procName  string       // 当前服务进程名称
	LinkRules []MatchRule  // 服务互联发现规则，用于控制哪些服务可以相互连接
)

// GetProcName 获取当前服务进程名称
// 返回:
//   - string: 进程名称，通常在Init函数中设置
func GetProcName() string {
	return procName
}

// GetWANIP 获取外网IP地址
// 外网IP用于通知客户端连接，例如login服务通知客户端game服务的外网IP
// 返回:
//   - string: 外网IP地址，如果未设置则返回空字符串
func GetWANIP() string {
	return flagWANIP
}

// GetSvcGroup 获取服务所在的分组
// 服务分组用于标识服务器所在的物理位置或逻辑分组
// 返回:
//   - string: 服务分组名称
func GetSvcGroup() string {
	return flagSvcGroup
}

// GetSvcIndex 获取服务索引
// 服务索引用于标识同类服务的不同进程实例，同类服务中的索引必须唯一
// 返回:
//   - string: 服务索引
func GetSvcIndex() string {
	return flagSvcIndex
}

// GetDiscoveryAddr 获取服务发现服务器的地址
// 返回:
//   - string: 服务发现服务器地址，格式为"host:port"
func GetDiscoveryAddr() string {
	return flagDiscoveryAddr
}

// GetCommtype 获取通信类型
// 返回:
//   - string: 通信类型标识
func GetCommtype() string {
	return flagCommType
}
