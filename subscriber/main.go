package main

import (
	"flag"
	"strconv"
	"time"

	config "go-perform-nats/config"
	util "go-perform-nats/util"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
)

func main() {
	// logrusの初期化
	logrus.SetLevel(logrus.ErrorLevel)

	// 引数
	arg_loop_ns := flag.Int("d", 5000*1000*1000, "duration ns")
	arg_iter := flag.Int("i", 20, "iteration count")
	args_subscripton_name := flag.String("n", "default", "subscription name")
	arg_thread := flag.Int("c", 1, "thread(gorouting) count")
	arg_topic_expr := flag.String("t", "test", "topic expression (e.g. /test/ or /test/:1,2,3 or /test/:1-3) ")
	arg_server_url := flag.String("s", config.Config_DefaultUrl, "connection string")
	arg_topic_divide := flag.Int("v", 1, "divide value")
	flag.Parse()

	//
	var iter = *arg_iter
	var thread = *arg_thread
	var loop_ns int64 = int64(*arg_loop_ns)

	// パフォーマンス測定用カウンタ生成
	pcMap := util.CreatePerformCounterMap()
	hist := util.CreateHistogram()

	// TopicSupplierFactory生成
	factory := util.CreateFactory().
		ParseTopicExpression(*arg_topic_expr).
		SetDistoribution(*arg_topic_divide)

	// Subscriberのセットアップ
	println("connecting")
	for i := 0; i < thread; i++ {
		thread_id := strconv.Itoa(i)
		supplier := factory.Build()
		msgChannel := make(chan pulsar.ConsumerMessage, 100)

		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL:               *arg_server_url,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second,
		})
		if err != nil {
			panic("Connect error:[" + err.Error() + "] thread:" + thread_id)
		}

		topics := supplier.GetAll()
		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topics:           topics,
			SubscriptionName: *args_subscripton_name + "<>" + thread_id,
			Type:             pulsar.Exclusive,
			MessageChannel:   msgChannel,
		})
		if err != nil {
			panic("Subscribe error:[" + err.Error() + "] thread:" + thread_id)
		}
		go func(msgChannel chan pulsar.ConsumerMessage) {
			for cm := range msgChannel {
				msg := cm.Message

				ping := util.DeserialPing(msg.Payload())
				key := thread_id + "<>" + msg.Topic()
				pcMap.Perform(key, &ping)
				hist.IncreamentPing(&ping)
			}
		}(msgChannel)

		defer consumer.Close()
	}

	println("connected")

	// パフォーマンス測定
	pcMap.CollectAndReset()
	var st int64 = time.Now().UnixNano()
	var et int64 = 0
	println("start-benchmark")
	for i := 0; i < iter; i++ {
		// 別スレッドで処理するため、メインスレッドはスリープする。
		time.Sleep(time.Duration(loop_ns) * time.Nanosecond)
		et = time.Now().UnixNano()

		snap := pcMap.CollectAndReset()
		snap.Print(et - st)
		st = et
	}
	hist.Print()
	println("end-benchmark")
}
