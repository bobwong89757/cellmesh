package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellmesh/util"
	"strings"
)

type MatchRule struct {
	Target string
}

func matchTarget(rule *MatchRule, desc *discovery.ServiceDesc) bool {

	return meshutil.WildcardPatternMatch(desc.GetMeta("SvcGroup"), rule.Target)
}
func ParseMatchRule(rule string) (ret []MatchRule) {

	for _, ruleStr := range strings.Split(rule, "|") {
		var rule MatchRule
		rule.Target = ruleStr
		ret = append(ret, rule)
	}

	return
}
