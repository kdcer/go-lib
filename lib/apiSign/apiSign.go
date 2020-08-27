package apiSign

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"strings"

	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

//接口签名参数
type SignParams struct {
	NormalSign      string   //参数标准签名盐值
	WebSign         string   //参数Web签名盐值   web易破解 设置独立盐值
	IgnoreFilterUrl []string //签名拦截忽略路径, 相对完整路径(不支持*)
	IgnoreParams    []string //忽略不参与签名的公共参数
	MasterKey       string   //万能密钥 直接跳过验证
}

//接口签名
//加密方式：MD5(盐值+|+排序普通参数[k=v&]+然后添加时间戳[k=v]+|+盐值) kv使用=和&连接，时间戳后面的&不用，盐值和参数之间用|隔开。
//比如：salt|appChannelId=10012&appProductCode=video_code&id=1&siteType=10003&unixTime=1583550103|salt
//设置开关-秘钥-参数排序-时间戳-时间是否过期-时区设置--web项目独立盐-上传等接口部分参数忽略-忽略部分接口
//header中有 1,appVersion 格式{appCode}-{appType}-{versionNumber} 取索引1为APP类型 (appType:1-Android, 2-Android(vue), 3-iOS,4-web) 平台
//			 2,unixTime Unix时间戳 3,sign 客户端加密结果
func APISign(r *ghttp.Request, signParams *SignParams) (err error) {
	// 静态页面不拦截
	if r.IsFileRequest() {
		return nil
	}
	//忽略路径直接跳过
	for _, signIgnorePath := range signParams.IgnoreFilterUrl {
		if signIgnorePath == r.URL.Path || strings.HasPrefix(r.URL.Path, signIgnorePath) {
			glog.Debug("忽略sign拦截链接url=", r.URL.String())
			return nil
		}
	}
	unixTime := r.Header.Get("unixTime")
	sign := r.Header.Get("sign")
	//万能签名直接跳过
	if sign == signParams.MasterKey {
		return nil
	}
	appVersion := r.Header.Get("appVersion")
	if unixTime == "" || sign == "" || appVersion == "" {
		return errors.New("参数不合法")
	}
	appVersionSlice := strings.Split(appVersion, "-")
	if len(appVersionSlice) < 3 {
		return errors.New("参数不合法")
	}
	platform := appVersionSlice[1]
	salt := signParams.NormalSign
	//如果是web 使用web独立的盐
	if platform == "4" {
		salt = signParams.WebSign
	}
	mp := r.GetRequestMap()

	//移除忽略参数
	for _, v := range signParams.IgnoreParams {
		delete(mp, v)
	}

	keys := make([]string, 0)
	for k := range mp {
		keys = append(keys, k)
	}
	//fmt.Println(keys)
	//排序
	sort.Strings(keys)
	//fmt.Println(keys)
	//拼接加密字符串
	kvString := bytes.NewBufferString(salt)
	kvString.WriteString("|")
	for _, v := range keys {
		kvString.WriteString(v)
		kvString.WriteString("=")
		kvString.WriteString(gconv.String(mp[v]))
		kvString.WriteString("&")
	}
	//参数之后添加时间
	kvString.WriteString("unixTime")
	kvString.WriteString("=")
	kvString.WriteString(unixTime)
	//最后添加盐
	kvString.WriteString("|")
	kvString.WriteString(salt)
	//fmt.Println(kvString.String())
	//fmt.Println(gurl.Decode(kvString.String()))// gf会自动decode 不用重复处理
	//生成服务器sign
	signServer, err := gmd5.Encrypt(kvString.String())
	if err != nil {
		glog.Error("参数签名 gmd5.Encrypt 错误", err)
		return errors.New("参数不合法")
	}
	//glog.Infof("服务器签名：%s", signServer)
	//签名不一致 直接返回
	if signServer != sign {
		glog.Infof("签名错误：MD5原始字符串=%s,服务端签名结果=%s,客户端签名=%s,appVersion=%s", kvString.String(), signServer, sign, appVersion)
		return errors.New("参数不合法")
	}
	//判断时间
	interval := gtime.Timestamp() - gconv.Int64(unixTime)
	//时间取绝对值 和服务器时间偏移60秒以上不合法
	interval = int64(math.Abs(float64(interval)))
	if interval > 60 {
		glog.Infof("签名超时：MD5原始字符串=%s,服务端签名结果=%s,客户端签名=%s,时间差=%d,appVersion=%s", kvString.String(), signServer, sign, interval, appVersion)
		return errors.New("参数不合法")
	}
	return nil
}
