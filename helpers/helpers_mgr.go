package helpers

import "github.com/bobwong89757/gnbutils/yaml"

var MConfig = &yaml.YamlUtil{
	KvsCache: make(map[string]map[string]string),
}