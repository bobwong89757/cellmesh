package meshutil

import (
	"flag"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/util"
)

// ApplyFlagFromFile 从文件中读取配置并应用到FlagSet
// 文件格式为键值对，每行一个配置项，格式为: key=value
// 参数:
//   - fs: flag.FlagSet实例
//   - filename: 配置文件路径
// 返回:
//   - error: 读取或应用失败时返回错误信息
func ApplyFlagFromFile(fs *flag.FlagSet, filename string) error {

	return util.ReadKVFile(filename, func(key, value string) bool {

		// 设置flagm
		fg := fs.Lookup(key)
		if fg != nil {
			log.GetLog().Infof("ApplyFlagFromFile: %s=%s", key, value)
			fg.Value.Set(value)
		} else {
			log.GetLog().Debugf("ApplyFlagFromFile: flag not found, %s", key)
		}

		return true
	})
}
