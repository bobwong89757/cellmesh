package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"sort"
)

// QueryServiceOp 是查询服务操作的枚举类型
// 用于控制查询过滤器的执行流程
type QueryServiceOp int

const (
	QueryServiceOp_NextFilter QueryServiceOp = iota // 继续下一个过滤器（内层循环）
	QueryServiceOp_NextDesc                         // 跳过当前服务，处理下一个服务（外层循环）
	QueryServiceOp_End                              // 终止所有遍历循环
)

// QueryResult 是查询结果的接口类型
// 用于在过滤器中返回查询结果
type QueryResult interface{}

// FilterFunc 是服务查询过滤函数类型
// 用于在查询服务时进行过滤和处理
// 返回值含义:
//   - true: 等效于QueryServiceOp_NextFilter，继续下一个过滤器
//   - false: 等效于QueryServiceOp_NextDesc，跳过当前服务
//   - QueryServiceOp_End: 终止所有遍历循环
//   - QueryResult类型: 作为查询结果返回，终止遍历
// 参数:
//   - *discovery.ServiceDesc: 当前处理的服务描述
// 返回:
//   - interface{}: 控制流程的值或查询结果
type FilterFunc func(*discovery.ServiceDesc) interface{}

// QueryService 根据服务名称查询服务，并通过过滤器链处理结果
// 参数:
//   - svcName: 要查询的服务名称
//   - filterList: 过滤器函数列表，按顺序执行
// 返回:
//   - ret: 查询结果，如果过滤器返回QueryResult类型则返回该值，否则返回nil
func QueryService(svcName string, filterList ...FilterFunc) (ret interface{}) {

	return QueryServiceEx(svcName, QueryServiceOption{}, filterList...)
}

// QueryServiceOption 是查询服务的选项配置
type QueryServiceOption struct {
	Sort bool // 是否对查询结果进行排序，按服务分组和索引排序
}

// QueryServiceEx 是QueryService的扩展版本，支持更多选项
// 参数:
//   - svcName: 要查询的服务名称
//   - opt: 查询选项
//   - filterList: 过滤器函数列表
// 返回:
//   - ret: 查询结果
func QueryServiceEx(svcName string, opt QueryServiceOption, filterList ...FilterFunc) (ret interface{}) {

	descList := discovery.Default.Query(svcName)

	if opt.Sort {
		sort.Slice(descList, func(i, j int) bool {

			a := descList[i]
			b := descList[j]

			aGroup := a.GetMeta("SvcGroup")
			bGroup := b.GetMeta("SvcGroup")

			if aGroup != bGroup {
				return aGroup < bGroup
			}

			aIndex := a.GetMeta("SvcIndex")
			bIndex := b.GetMeta("SvcIndex")

			return aIndex < bIndex
		})
	}

	for _, desc := range descList {

		for _, filter := range filterList {

			if filter == nil {
				continue
			}

			op := filter(desc)

			switch raw := op.(type) {
			case QueryServiceOp:
				switch raw {
				case QueryServiceOp_NextFilter:
				case QueryServiceOp_NextDesc:
					goto NextDesc
				case QueryServiceOp_End:
					return
				}
			case bool:
				if !raw {
					goto NextDesc
				}
			case QueryResult:
				ret = raw
			default:
				panic("unknown filter result")
			}
		}

	NextDesc:
	}

	return
}

// Filter_MatchSvcGroup 创建一个匹配指定服务组的过滤器
// 如果服务组为空字符串，则匹配所有服务
// 参数:
//   - svcGroup: 要匹配的服务组名称，空字符串表示匹配所有
// 返回:
//   - FilterFunc: 过滤器函数
func Filter_MatchSvcGroup(svcGroup string) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		if svcGroup == "" {
			return true
		}

		return desc.GetMeta("SvcGroup") == svcGroup
	}
}

// Filter_MatchSvcID 创建一个匹配指定服务ID的过滤器
// 如果找到匹配的服务，会返回该服务的描述信息作为查询结果
// 参数:
//   - svcid: 要匹配的服务ID
// 返回:
//   - FilterFunc: 过滤器函数
func Filter_MatchSvcID(svcid string) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		if desc.ID == svcid {
			return QueryResult(desc)
		}

		return true
	}
}

// Filter_MatchRule 创建一个匹配指定规则的过滤器
// 使用通配符模式匹配服务组，任意规则满足即可通过
// 参数:
//   - rules: 匹配规则列表
// 返回:
//   - FilterFunc: 过滤器函数
func Filter_MatchRule(rules []MatchRule) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		// 任意规则满足即可
		for _, rule := range rules {

			if matchTarget(&rule, desc) {
				return true
			}
		}

		return false
	}

}
