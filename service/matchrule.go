package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellmesh/util"
	"strings"
)

// MatchRule 是服务匹配规则
// 用于控制哪些服务可以相互连接
type MatchRule struct {
	Target string // 目标匹配模式，支持通配符，用于匹配服务组
}

// matchTarget 检查服务描述是否匹配规则
// 使用通配符模式匹配服务组
// 参数:
//   - rule: 匹配规则
//   - desc: 服务描述
// 返回:
//   - bool: 如果匹配返回true
func matchTarget(rule *MatchRule, desc *discovery.ServiceDesc) bool {

	return meshutil.WildcardPatternMatch(desc.GetMeta("SvcGroup"), rule.Target)
}
// ParseMatchRule 解析匹配规则字符串
// 规则字符串使用"|"分隔多个规则，例如: "group1|group2"
// 参数:
//   - rule: 规则字符串
// 返回:
//   - ret: 解析后的规则列表
func ParseMatchRule(rule string) (ret []MatchRule) {

	for _, ruleStr := range strings.Split(rule, "|") {
		var rule MatchRule
		rule.Target = ruleStr
		ret = append(ret, rule)
	}

	return
}
