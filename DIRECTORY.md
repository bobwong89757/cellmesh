# Cellmesh 项目目录结构

本文档详细说明了cellmesh项目的目录结构和每个文件的作用。

## 项目根目录

```
cellmesh/
├── discovery/          # 服务发现相关代码
├── service/            # 服务通信基础代码
├── util/               # 通用工具代码
├── helpers/            # 辅助工具
├── tool/               # 代码生成工具
├── dummy.go            # 占位文件，用于go get
├── go.mod              # Go模块依赖管理文件
├── go.sum              # Go模块依赖校验文件
├── LICENSE             # 许可证文件
├── README.md           # 项目说明文档
└── DIRECTORY.md        # 本文件，目录结构说明文档
```

## 详细目录结构

### 根目录文件

- **dummy.go**: 占位文件，用于让`go get`能够正确识别和下载此包
- **go.mod**: Go模块依赖管理文件，定义项目依赖
- **go.sum**: Go模块依赖校验文件，包含依赖包的校验和
- **LICENSE**: 项目许可证文件（MIT License）
- **README.md**: 项目主要说明文档，包含使用说明、特性介绍等

---

## discovery/ - 服务发现包

服务发现相关的核心代码，提供统一的服务发现接口和实现。

```
discovery/
├── discovery.go        # 服务发现接口定义
├── desc.go             # 服务描述结构体定义
├── safevalue.go        # 大值安全存储和读取
├── safevalue_test.go   # safevalue的单元测试
├── util.go             # 服务发现工具函数
├── kvconfig/           # KV配置快速获取接口
│   └── kvconfig.go     # 配置获取辅助函数
└── memsd/              # memsd服务发现实现
    ├── api/            # memsd客户端API
    ├── deps/           # memsd服务器依赖
    ├── model/          # memsd数据模型
    ├── proto/          # memsd协议定义
    └── memsd           # memsd服务器可执行文件
```

### discovery/ 文件说明

- **discovery.go**: 
  - 定义`Discovery`接口，提供统一的服务发现抽象
  - 定义`ValueMeta`结构体，用于KV配置元数据
  - 定义`CheckerFunc`健康检查函数类型
  - 提供`Default`全局服务发现实例

- **desc.go**: 
  - 定义`ServiceDesc`服务描述结构体
  - 提供服务描述的比较、元数据操作、地址格式化等方法
  - 包含服务描述的字符串表示方法

- **safevalue.go**: 
  - 提供大值分片存储功能（超过300KB的值自动分片）
  - 支持值的压缩和解压缩
  - 实现`SafeSetValue`和`SafeGetValue`函数

- **safevalue_test.go**: 
  - `safevalue.go`的单元测试文件

- **util.go**: 
  - `BytesToAny`: 字节数组到任意类型的转换
  - `AnyToBytes`: 任意类型到字节数组的转换
  - `ValueMetaToSlice`: ValueMeta数组到切片的转换

### discovery/kvconfig/ - KV配置包

- **kvconfig.go**: 
  - 提供类型安全的配置获取函数：`String`、`Int32`、`Int64`、`Bool`
  - 配置不存在时自动使用默认值并写入服务发现系统

### discovery/memsd/ - memsd服务发现实现

memsd是cellmesh自研的轻量级服务发现系统，面向游戏服务优化。

#### discovery/memsd/api/ - 客户端API

```
api/
├── api.go          # memsd客户端主入口和核心实现
├── api_test.go     # API单元测试
├── config.go       # 客户端配置结构
├── conn.go         # 连接管理
├── const.go        # 常量定义
├── kv.go           # KV操作实现
├── packet.go       # 数据包处理
├── rpc.go          # RPC调用实现
├── setup.go        # 初始化设置
├── svc.go          # 服务注册和查询实现
├── transmitter.go  # 消息传输器
└── util.go         # API工具函数
```

- **api.go**: 
  - `memDiscovery`结构体，实现`Discovery`接口
  - `NewDiscovery`函数，创建memsd客户端实例
  - 管理连接、缓存、通知等核心逻辑

- **config.go**: 
  - `Config`配置结构体
  - `DefaultConfig`默认配置函数

- **conn.go**: 
  - 连接建立和管理
  - 处理连接事件和消息

- **kv.go**: 
  - KV配置的增删改查实现
  - 缓存管理

- **svc.go**: 
  - 服务注册和注销
  - 服务查询和缓存更新
  - 服务变化通知

- **rpc.go**: 
  - 远程过程调用实现
  - 请求响应处理

- **packet.go**: 
  - 数据包编解码

- **transmitter.go**: 
  - 消息传输器实现

- **setup.go**: 
  - 初始化设置

- **util.go**: 
  - API相关的工具函数

#### discovery/memsd/deps/ - 服务器依赖

```
deps/
├── cmd.go          # 命令行工具实现
├── persist.go      # 数据持久化
├── redundant.go    # 冗余处理
├── sd.go           # 服务发现服务器初始化
├── svc.go          # 服务处理
└── svc_msg.go      # 服务消息处理
```

- **sd.go**: 
  - `InitSD`函数，初始化服务发现客户端
  - `DiscoveryExtend`扩展接口定义

- **svc.go**: 
  - `StartSvc`函数，启动memsd服务器
  - 服务消息处理

- **svc_msg.go**: 
  - 服务相关消息的处理逻辑

- **cmd.go**: 
  - 命令行工具实现（查看服务、查看配置、设置值等）

- **persist.go**: 
  - 数据持久化到文件
  - 从文件加载数据

- **redundant.go**: 
  - 冗余处理逻辑

#### discovery/memsd/model/ - 数据模型

```
model/
├── kv.go           # KV存储模型
└── svcmodel.go     # 服务模型
```

- **kv.go**: 
  - `ValueMeta`结构体，KV存储的元数据
  - `PersistFile`持久化文件结构
  - KV的增删改查操作

- **svcmodel.go**: 
  - 服务相关的模型定义
  - 服务键前缀、UUID生成器等

#### discovery/memsd/proto/ - 协议定义

```
proto/
├── sd.proto            # Protocol Buffers协议定义
├── msgbind_gen.go      # 自动生成的消息绑定代码
├── msgsvc_gen.go       # 自动生成的消息服务代码
├── proto_svc.txt       # 协议服务列表
├── MakeProto.sh        # 协议生成脚本
└── protolist.sh        # 协议列表脚本
```

- **sd.proto**: 
  - memsd的Protocol Buffers协议定义
  - 定义所有消息类型

- **msgbind_gen.go**: 
  - 自动生成的消息绑定代码

- **msgsvc_gen.go**: 
  - 自动生成的消息服务处理代码

- **MakeProto.sh**: 
  - 生成Protocol Buffers代码的脚本

- **protolist.sh**: 
  - 列出所有协议的脚本

- **proto_svc.txt**: 
  - 协议服务列表文件

#### discovery/memsd/memsd

- **memsd**: 
  - memsd服务器的可执行文件（编译后的二进制）

---

## service/ - 服务通信包

服务通信基础代码，提供服务的注册、发现、连接管理等功能。

```
service/
├── discovery.go        # 服务发现和连接
├── flag.go             # 命令行参数定义
├── hooker.go           # 服务互联消息处理Hooker
├── init.go             # 服务初始化
├── matchrule.go        # 服务匹配规则
├── model.go            # 服务模型和全局变量
├── msg.go              # 服务消息定义
├── multipeer.go        # 多Peer管理
├── query.go            # 服务查询和过滤
├── reg.go              # 服务注册
├── remotesvc.go        # 远程服务管理
├── safevalue_test.go   # safevalue测试
└── svcid_test.go       # svcid测试
└── svcid.go            # 服务ID生成和解析
```

### service/ 文件说明

- **discovery.go**: 
  - `DiscoveryService`函数，发现并连接到指定服务
  - `DiscoveryOption`服务发现选项配置

- **flag.go**: 
  - 定义服务相关的命令行参数变量
  - `InitServerConfig`初始化服务器配置

- **hooker.go**: 
  - `SvcEventHooker`服务互联消息处理Hooker
  - 处理服务间的连接建立、身份确认等事件
  - 注册`tcp.svc`和`tcp.client`处理器

- **init.go**: 
  - `Init`初始化服务框架
  - `ConnectDiscovery`连接到服务发现服务器
  - `LogParameter`打印服务参数
  - `WaitExitSignal`等待退出信号

- **matchrule.go**: 
  - `MatchRule`服务匹配规则结构体
  - `ParseMatchRule`解析匹配规则字符串
  - 支持通配符模式匹配

- **model.go**: 
  - 定义全局变量：`procName`、`LinkRules`
  - 提供获取服务参数的函数：`GetProcName`、`GetWANIP`、`GetSvcGroup`等

- **msg.go**: 
  - `ServiceIdentifyACK`服务身份确认消息
  - `GetPassThrough`从relay事件提取透传数据
  - `Reply`回复消息的便捷函数

- **multipeer.go**: 
  - `MultiPeer`接口，管理多个Peer连接
  - `multiPeer`实现，用于连接多个服务实例

- **query.go**: 
  - `QueryService`查询服务并应用过滤器
  - `FilterFunc`过滤器函数类型
  - 提供多种过滤器：`Filter_MatchSvcGroup`、`Filter_MatchSvcID`、`Filter_MatchRule`

- **reg.go**: 
  - `Register`将Acceptor注册到服务发现系统
  - `Unregister`注销服务
  - `ServiceMeta`服务元数据类型

- **remotesvc.go**: 
  - `RemoteServiceContext`远程服务上下文
  - `AddRemoteService`添加远程服务
  - `RemoveRemoteService`移除远程服务
  - `GetRemoteService`根据服务ID获取会话
  - `VisitRemoteService`遍历远程服务

- **svcid.go**: 
  - `MakeSvcID`构造服务ID
  - `MakeLocalSvcID`构造本地服务ID
  - `GetLocalSvcID`获取本进程服务ID
  - `ParseSvcID`解析服务ID

- **svcid_test.go**: 
  - `svcid.go`的单元测试

- **safevalue_test.go**: 
  - safevalue的测试文件

---

## util/ - 工具包

通用工具代码，提供UUID生成、通配符匹配、配置文件读取等功能。

```
util/
├── uuid64.go           # 64位UUID生成器
├── uuid64_test.go      # UUID生成器测试
├── wilecard.go         # 通配符模式匹配
├── wilecard_test.go    # 通配符匹配测试
└── flagfile.go         # 从文件读取Flag配置
```

### util/ 文件说明

- **uuid64.go**: 
  - `UUID64Generator`64位UUID生成器
  - `UUID64Component`UUID组件
  - 支持时间戳、序列号、固定值等组件类型
  - `AddTimeComponent`、`AddSeqComponent`、`AddConstComponent`等方法

- **uuid64_test.go**: 
  - UUID生成器的单元测试

- **wilecard.go**: 
  - `WildcardPatternMatch`通配符模式匹配函数
  - 支持`?`（单个字符）和`*`（多个字符）通配符
  - 使用动态规划算法实现

- **wilecard_test.go**: 
  - 通配符匹配的单元测试

- **flagfile.go**: 
  - `ApplyFlagFromFile`从文件读取配置并应用到FlagSet
  - 支持键值对格式的配置文件

---

## helpers/ - 辅助工具包

辅助工具代码。

```
helpers/
└── helpers_mgr.go      # 辅助工具管理器
```

### helpers/ 文件说明

- **helpers_mgr.go**: 
  - `MConfig`全局YAML配置工具实例
  - 提供YAML配置的读取和缓存功能

---

## tool/ - 工具包

代码生成工具，用于生成协议相关的代码。

```
tool/
└── protogen/           # 协议生成器
    ├── main.go         # 主程序入口
    └── gengo/          # Go代码生成
        ├── func.go     # 函数生成
        ├── gen.go      # 代码生成核心
        └── text.go     # 文本模板处理
```

### tool/protogen/ 文件说明

- **main.go**: 
  - 协议生成器的主程序入口
  - 解析命令行参数并调用生成器

- **gengo/func.go**: 
  - 函数生成的辅助函数

- **gengo/gen.go**: 
  - 代码生成的核心逻辑

- **gengo/text.go**: 
  - 文本模板处理相关函数

---

## 总结

cellmesh项目采用清晰的模块化设计：

1. **discovery包**: 提供统一的服务发现接口和memsd实现
2. **service包**: 提供服务通信的基础功能
3. **util包**: 提供通用的工具函数
4. **helpers包**: 提供辅助工具
5. **tool包**: 提供代码生成工具

每个包都有明确的职责，便于维护和扩展。

