package meshutil

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// numFToMask 将半字节数量转换为对应的位掩码
// numF表示半字节数量，1个半字节=4位=1个十六进制位
// 参数:
//   - numF: 半字节数量，范围1-16
// 返回:
//   - uint64: 对应的位掩码
func numFToMask(numF uint) uint64 {
	switch numF {
	case 1:
		return 0xF
	case 2:
		return 0xFF
	case 3:
		return 0xFFF
	case 4:
		return 0xFFFF
	case 5:
		return 0xFFFFF
	case 6:
		return 0xFFFFFF
	case 7:
		return 0xFFFFFFF
	case 8:
		return 0xFFFFFFFF
	case 9:
		return 0xFFFFFFFFF
	case 10:
		return 0xFFFFFFFFFF
	case 11:
		return 0xFFFFFFFFFFF
	case 12:
		return 0xFFFFFFFFFFFF
	case 13:
		return 0xFFFFFFFFFFFFF
	case 14:
		return 0xFFFFFFFFFFFFFF
	case 15:
		return 0xFFFFFFFFFFFFFFF
	case 16:
		return 0xFFFFFFFFFFFFFFFF
	default:
		panic("numF shound in range 1~16")
	}
}

// UUID64Component 是UUID64生成器的组件
// 每个组件负责生成UUID的一部分
type UUID64Component struct {
	ValueSrc func() uint64 // 数值来源函数，每次调用返回一个数值
	NumF     uint          // 占用的半字节数量，1个F=0.5个字节=4位
}

// UUID64Generator 是64位UUID生成器
// 通过组合多个组件来生成唯一的64位ID
// 支持时间戳、序列号、固定值等组件类型
type UUID64Generator struct {
	seqGen   uint64              // 序列号生成器的当前值
	comSet   []*UUID64Component  // 组件列表
	genGuard sync.Mutex          // 保护生成过程的互斥锁
}

const (
	MaxNumFInt64 = 16 // int64类型最多支持16个半字节（64位）
)

// AddComponent 添加一个组件到生成器
// 参数:
//   - com: 要添加的组件
// 注意: 所有组件的总半字节数不能超过16，否则会panic
func (self *UUID64Generator) AddComponent(com *UUID64Component) {

	// 检查范围
	numFToMask(com.NumF)

	self.comSet = append(self.comSet, com)

	// 检查总组件超过位数
	if self.UsedNumF() > MaxNumFInt64 {
		panic("total bit over int64(8bit) range")
	}
}

// UsedNumF 计算已使用的半字节总数
// 返回:
//   - ret: 已使用的半字节数量
func (self *UUID64Generator) UsedNumF() (ret uint) {

	for _, com := range self.comSet {
		ret += com.NumF
	}

	return
}

// LeftNumF 计算剩余的可用半字节数
// 返回:
//   - ret: 剩余的半字节数量
func (self *UUID64Generator) LeftNumF() (ret uint) {
	return MaxNumFInt64 - self.UsedNumF()
}

// AddSeqComponent 添加一个序列号组件
// 序列号组件每次生成时自动递增
// 参数:
//   - numF: 占用的半字节数量
//   - init: 初始序列号值
// 注意: 初始值必须在指定范围内，否则会panic
func (self *UUID64Generator) AddSeqComponent(numF uint, init uint64) {

	if init&numFToMask(numF) != init {
		panic(fmt.Sprintf("const component out of range, expect 0~%d, got, %d", numFRange(numF), init))
	}

	self.seqGen = init
	self.AddComponent(&UUID64Component{
		ValueSrc: func() uint64 {
			self.seqGen++
			return self.seqGen
		},

		NumF: numF,
	})

}

// numFRange 计算指定半字节数量能表示的最大值
// 参数:
//   - numF: 半字节数量
// 返回:
//   - uint: 最大值（16^numF - 1）
func numFRange(numF uint) uint {

	return uint(math.Pow(16, float64(numF))) - 1
}

// AddConstComponent 添加一个固定值组件
// 固定值组件每次生成时返回相同的值
// 参数:
//   - numF: 占用的半字节数量
//   - constNumber: 固定值
// 注意: 固定值必须在指定范围内，否则会panic
func (self *UUID64Generator) AddConstComponent(numF uint, constNumber uint64) {

	uconst := uint64(constNumber)

	if uconst&numFToMask(numF) != constNumber {
		panic(fmt.Sprintf("const component out of range, expect 0~%d, got, %d", numFRange(numF), constNumber))
	}

	self.AddComponent(&UUID64Component{
		ValueSrc: func() uint64 {
			return uconst
		},

		NumF: numF,
	})

}

const timeStartPoint = 946656000 // 时间参考点：2000/1/1 0:0:0的Unix时间戳，用于延迟2039年时间戳溢出问题

// AddTimeComponent 添加一个时间戳组件
// 时间戳组件返回当前时间相对于参考点的秒数
// 参数:
//   - numF: 占用的半字节数量
func (self *UUID64Generator) AddTimeComponent(numF uint) {
	self.AddComponent(&UUID64Component{
		ValueSrc: func() uint64 {
			return uint64(time.Now().Unix() - timeStartPoint)
		},

		NumF: numF,
	})

}

// Generate 按照组件规则生成一个64位UUID
// 从右到左组合各个组件的值
// 返回:
//   - ret: 生成的64位UUID
func (self *UUID64Generator) Generate() (ret uint64) {

	self.genGuard.Lock()
	var offset uint
	for i := len(self.comSet) - 1; i >= 0; i-- {
		g := self.comSet[i]
		mask := numFToMask(g.NumF)
		part := (g.ValueSrc() & mask) << offset
		ret |= part
		offset += g.NumF * 4
	}

	self.genGuard.Unlock()

	return
}

// NewUUID64Generator 创建一个新的UUID64生成器
// 返回:
//   - *UUID64Generator: UUID生成器实例
func NewUUID64Generator() *UUID64Generator {

	return &UUID64Generator{}
}
