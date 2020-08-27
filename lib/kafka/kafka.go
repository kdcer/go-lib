package kafka

import (
	"strings"
	"time"

	"github.com/gogf/gf/os/glog"

	"github.com/Shopify/sarama"
	"github.com/gogf/gf/frame/g"
)

var address = strings.Split(g.Cfg().GetString("kafka.host"), ",")
var Producer sarama.AsyncProducer

func Init() {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1

	glog.Info("start make producer")
	//使用配置,新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(address, config)
	if e != nil {
		glog.Errorf("kafka创建生产者失败，err=", e)
		return
	}
	Producer = producer
	//defer producer.AsyncClose()
}

func Input(data *Message) {
	// 发送的消息,主题。
	msg := &sarama.ProducerMessage{
		Topic:     data.Topic,
		Value:     sarama.ByteEncoder(data.Body), //将字节数组转换为Encoder
		Timestamp: time.Now(),
	}
	SaramaProducer(msg)
}

//异步消息
func SaramaProducer(msg *sarama.ProducerMessage) {
	//判断哪个通道发送过来数据.
	go func(p sarama.AsyncProducer) {
		defer func() {
			if err := recover(); err != nil {
				glog.Error("kafka连接异常", err)
			}
		}()
		select {
		case suc := <-p.Successes():
			//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			glog.Info("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
		case fail := <-p.Errors():
			//fmt.Println("err: ", fail.Err)
			glog.Error("err: ", fail.Err)
		}
	}(Producer)
	//使用通道发送
	Producer.Input() <- msg
}
