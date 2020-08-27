package util

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
	"io"
	"math/rand"

	"github.com/gogf/gf/util/gconv"
)

//生成32位md5字串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成uuid字串  使用此方法获取随机Guid
func UUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

//根据可变string拼接字符串
func NewBufferString(s ...string) string {
	l := len(s)
	if l == 0 {
		return ""
	}
	if l == 1 {
		return s[0]
	}
	b := bytes.NewBufferString(s[0])
	for i := 1; i < l; i++ {
		b.WriteString(s[i])
	}
	return b.String()
}

// String hashes a string to a unique hashcode.
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
// 并不能保证是绝对唯一的 比如：'49d3f976734b089fe0ba28960948459d' 和 'f0a8bbfdcf08874262abb3aeabdeaa69' 这两个字符串返回的结果是一样的都是687199550  只能用于不那么严格的场景
func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	} else {
		return -v
	}
}

// 生成随机数据查询sql，先查询总数count，随机获取limit个不重复的offset(从0到count)，将这些随机数作为limit的偏移量 ，例如limit 10,1,limit 30,1。sql语句用union all 拼接。
// count 有效数据总数
// limit 查询的总数
// tableName 表名
// fields 字段
// conditions 条件
// order 排序
// return 生成的sql
func GenerateRandDataSql(count, limit int, tableName, fields, conditions, order string) string {
	offsets := GetRandIds(count, limit)
	sqlStr := bytes.NewBufferString("")
	lens := len(offsets)
	for i, v := range offsets {
		sqlStr.WriteString("(SELECT ")
		sqlStr.WriteString(fields)
		sqlStr.WriteString(" FROM ")
		sqlStr.WriteString(tableName)
		sqlStr.WriteString(" WHERE ")
		sqlStr.WriteString(conditions)
		sqlStr.WriteString(order)
		sqlStr.WriteString(" LIMIT " + gconv.String(v) + ",1) \n")
		if i < lens-1 {
			sqlStr.WriteString("union all \n")
		}
	}
	sqlStr.WriteString(";")
	return sqlStr.String()
}

// 随机查询数据，根据数据总量和查询数量返回随机数offset切片
// count 有效数据总数
// limit 查询的总数
// return 偏移量切片
func GetRandIds(count, limit int) []int {
	// 如果查询数量大于总数，只查询count个
	if limit > count {
		limit = count
	}
	rand.Seed(int64(HashCode(UUID()))) // 使用hashcode设置随机种子
	idMap := map[int]byte{}            // 用于判断是否已经存储
	offsets := []int{}                 // 用于存储的切片
	// 使用死循环，获取足够limit的随机数
	for {
		index := rand.Intn(count) // 从1到count-1
		// 如果没有存储才进行存储
		if _, ok := idMap[index]; !ok {
			idMap[index] = 0
			offsets = append(offsets, index)
		}
		// 获取足够的数据则返回
		if len(offsets) == limit {
			break
		}
	}
	return offsets
}

// 通过两重循环过滤重复元素  时间换空间
func RemoveRepByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 通过map主键唯一的特性过滤重复元素 空间换时间
func RemoveRepByMap(slc []string) []string {
	slcLen := len(slc)
	result := make([]string, 0, slcLen)
	tempMap := make(map[string]struct{}, slcLen) // 存放不重复主键
	var l = 0
	for _, e := range slc {
		l = len(tempMap)
		tempMap[e] = struct{}{}
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// 字符截取指定长度(str原字符串，rep替换字符串,limit长度限制)
func StringsTruncate(str, rep string, limit int) (dst string) {
	if len([]rune(str)) > limit {
		dst = string([]rune(str)[:limit]) + rep
	} else {
		dst = str
	}
	return
}
