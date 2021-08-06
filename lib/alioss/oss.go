package alioss

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"

	"github.com/gogf/gf/os/glog"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

// OSS .
type OSS struct {
	client *oss.Client
	bucket *oss.Bucket
	*OssConfig
}

var (
	syncOnce sync.Once
	_oss     *OSS
)

// New .
func New(ossConfig *OssConfig) *OSS {
	syncOnce.Do(func() {
		_oss = new(OSS)
		_oss.OssConfig = ossConfig
		client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		_oss.client = client
		bucket, err := client.Bucket(ossConfig.BucketName)
		if err != nil {
			glog.Error(err)
			panic(err)
		}
		_oss.bucket = bucket
	})
	return _oss
}

// Put 上传.
func (oss *OSS) Put(objectKey string, reader io.Reader, fileName string) (uri string, err error) {
	err = oss.bucket.PutObject(objectKey, reader)
	if err != nil {
		glog.Error(err)
	}
	//uri = fmt.Sprintf("https://%s/%s/%s", oss.bucket.GetConfig().Endpoint, time.Now().Format("2006-01-02"), objectKey)
	// 文件名编码返回url
	objectKey = strings.Replace(objectKey, fileName, url.PathEscape(fileName), 1)
	uri = fmt.Sprintf("https://%s.%s/%s", oss.BucketName, strings.Replace(oss.Endpoint, "https://", "", 1), objectKey)
	return
}

// Delete 删除.
func (oss *OSS) Delete(objectKey string) (err error) {
	err = oss.bucket.DeleteObject(objectKey)
	if err != nil {
		glog.Error(err)
	}
	return
}

// ReName 重命名.
func (oss *OSS) ReName(srcObjectKey, destObjectKey string) (err error) {
	// 移除url前缀获取objectName
	urlPrefix := fmt.Sprintf("https://%s.%s/", oss.BucketName, strings.Replace(oss.Endpoint, "https://", "", 1))
	srcObjectKey = strings.Replace(srcObjectKey, urlPrefix, "", 1)
	destObjectKey = strings.Replace(destObjectKey, urlPrefix, "", 1)
	_, err = oss.bucket.CopyObject(srcObjectKey, destObjectKey)
	if err != nil {
		glog.Error(err)
		return
	}
	err = oss.bucket.DeleteObject(srcObjectKey)
	if err != nil {
		glog.Error(err)
	}
	return
}
