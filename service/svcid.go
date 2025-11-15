package service

import (
	"errors"
	"fmt"
	"strconv"
)

// 服务ID格式说明: svcName#svcIndex@svcGroup
// 例如: "game#1@group1" 表示group1分组中的第1个game服务
// 这种格式保证了服务ID的全局唯一性

// MakeSvcID 构造服务ID
// 参数:
//   - svcName: 服务名称
//   - svcIndex: 服务索引
//   - svcGroup: 服务分组
// 返回:
//   - string: 格式化的服务ID
func MakeSvcID(svcName string, svcIndex int, svcGroup string) string {
	return fmt.Sprintf("%s#%d@%s", svcName, svcIndex, svcGroup)
}

// MakeLocalSvcID 构造本地服务的ID
// 使用当前进程的服务索引和分组信息
// 参数:
//   - svcName: 服务名称
// 返回:
//   - string: 格式化的服务ID
func MakeLocalSvcID(svcName string) string {
	index,_ := strconv.Atoi(flagSvcIndex)
	return MakeSvcID(svcName,index, flagSvcGroup)
}

// GetLocalSvcID 获取本进程的服务ID
// 使用当前进程名称、索引和分组信息构造
// 返回:
//   - string: 当前进程的服务ID
func GetLocalSvcID() string {
	return MakeLocalSvcID(GetProcName())
}

// ParseSvcID 解析服务ID，提取服务名称、索引和分组
// 服务ID格式: svcName#svcIndex@svcGroup
// 参数:
//   - svcid: 要解析的服务ID字符串
// 返回:
//   - svcName: 服务名称
//   - svcIndex: 服务索引
//   - svcGroup: 服务分组
//   - err: 解析失败时返回错误信息
func ParseSvcID(svcid string) (svcName string, svcIndex int, svcGroup string, err error) {

	var sharpPos, atPos = -1, -1

	for pos, c := range svcid {
		switch c {
		case '#':
			sharpPos = pos
			svcName = svcid[:sharpPos]
		case '@':
			atPos = pos
			svcGroup = svcid[atPos+1:]

			if sharpPos == -1 {
				break
			}

			var n int64
			n, err = strconv.ParseInt(svcid[sharpPos+1:atPos], 10, 32)
			if err != nil {
				break
			}
			svcIndex = int(n)
		}
	}

	if sharpPos == -1 {
		err = errors.New("missing '#' in svcid")
	}

	if atPos == -1 {
		err = errors.New("missing '@' in svcid")
	}

	return
}
