package meshutil

// stringToRuneSlice 将字符串转换为rune切片
// 用于支持Unicode字符的通配符匹配
// 参考: https://siongui.github.io/2017/04/11/go-wildcard-pattern-matching/
// 参数:
//   - s: 源字符串
// 返回:
//   - []rune: rune切片
func stringToRuneSlice(s string) []rune {
	var r []rune
	for _, runeValue := range s {
		r = append(r, runeValue)
	}
	return r
}

// initLookupTable 初始化动态规划查找表
// 用于通配符模式匹配的缓存表
// 参数:
//   - row: 行数
//   - column: 列数
// 返回:
//   - [][]bool: 初始化的二维布尔数组
func initLookupTable(row, column int) [][]bool {
	lookup := make([][]bool, row)
	for i := range lookup {
		lookup[i] = make([]bool, column)
	}
	return lookup
}

// WildcardPatternMatch 使用通配符模式匹配字符串
// 支持的通配符:
//   - '?' 匹配任意单个字符
//   - '*' 匹配任意多个字符（包括0个）
// 使用动态规划算法实现
// 参数:
//   - str: 要匹配的字符串
//   - pattern: 通配符模式
// 返回:
//   - bool: 如果匹配返回true，否则返回false
func WildcardPatternMatch(str, pattern string) bool {
	s := stringToRuneSlice(str)
	p := stringToRuneSlice(pattern)

	// empty pattern can only match with empty string
	if len(p) == 0 {
		return len(s) == 0
	}

	// lookup table for storing results of subproblems
	// zero value of lookup is false
	lookup := initLookupTable(len(s)+1, len(p)+1)

	// empty pattern can match with empty string
	lookup[0][0] = true

	// Only '*' can match with empty string
	for j := 1; j < len(p)+1; j++ {
		if p[j-1] == '*' {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	// fill the table in bottom-up fashion
	for i := 1; i < len(s)+1; i++ {
		for j := 1; j < len(p)+1; j++ {
			if p[j-1] == '*' {
				// Two cases if we see a '*'
				// a) We ignore ‘*’ character and move
				//    to next  character in the pattern,
				//     i.e., ‘*’ indicates an empty sequence.
				// b) '*' character matches with ith
				//     character in input
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]

			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				// Current characters are considered as
				// matching in two cases
				// (a) current character of pattern is '?'
				// (b) characters actually match
				lookup[i][j] = lookup[i-1][j-1]

			} else {
				// If characters don't match
				lookup[i][j] = false
			}
		}
	}

	return lookup[len(s)][len(p)]
}
