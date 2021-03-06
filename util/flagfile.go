package meshutil

import (
	"flag"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/util"
)

func ApplyFlagFromFile(fs *flag.FlagSet, filename string) error {

	return util.ReadKVFile(filename, func(key, value string) bool {

		// 设置flagm
		fg := fs.Lookup(key)
		if fg != nil {
			log.GetLog().Info("ApplyFlagFromFile: %s=%s", key, value)
			fg.Value.Set(value)
		} else {
			log.GetLog().Error("ApplyFlagFromFile: flag not found, %s", key)
		}

		return true
	})
}
