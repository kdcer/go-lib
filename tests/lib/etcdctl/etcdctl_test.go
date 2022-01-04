package etcdctl

import (
	"context"
	"fmt"
	"testing"

	"go.etcd.io/etcd/clientv3"

	"github.com/gogf/gf/os/glog"

	"github.com/kdcer/go-lib/lib/etcdctl"
)

func Test_Lock(t *testing.T) {
	etcd := etcdctl.New([]string{"127.0.0.1:2379"})
	res := etcd.Lock("lockKey", func() interface{} {
		fmt.Println(1)
		return 1
	})

	fmt.Println(res)
	_, err := etcd.Put(context.Background(), "go/test/a", "1")
	if err != nil {
		glog.Error(err)
	}

	_, err = etcd.Get(context.Background(), "go/test/a", clientv3.WithPrefix())
	if err != nil {
		glog.Error(err)
	}

	resp, err := etcd.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		glog.Error(err)
	}
	fmt.Println(resp)
}
