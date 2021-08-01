package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"go-perform-nats/config"
	util "go-perform-nats/util"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.ErrorLevel)
	// 引数の処理
	arg_loop_ns := flag.Int("d", 5000*1000*1000, "duration ns")
	arg_iter := flag.Int("i", 20, "iteration count")
	arg_msg_size := flag.Int("m", 150, "message size(bytes)")
	arg_topic_expr := flag.String("t", "test", "topic expression")
	arg_server_url := flag.String("s", config.Config_DefaultUrl, "connection string")
	flag.Parse()

	//
	var loop_ns int64 = int64(*arg_loop_ns)
	var msg_size int = *arg_msg_size

	//
	factory := util.CreateFactory().ParseTopicExpression(*arg_topic_expr)

	// Connect to a server
	supplier := factory.Build()
	println("connecting")
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               *arg_server_url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		panic("connect error:" + err.Error())
	}
	defer client.Close()
	println("connected")

	topics := supplier.GetAll()
	producers := make(map[string]pulsar.Producer)
	for _, topic := range topics {
		producer, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: topic,
		})
		if err != nil {
			fmt.Errorf("Could not instantiate producer: %v", err)
		}
		producers[topic] = producer
		defer producer.Close()
	}

	// パフォーマンス測定開始
	println("start-benchmark")
	seqMap := make(map[string]uint32)
	for i := 0; i < *arg_iter; i++ {
		var cnt int64 = 0
		var st int64 = time.Now().UnixNano()
		var et int64 = 0

		for et-st < loop_ns {
			et = time.Now().UnixNano()
			topic := supplier.Get()
			seq := seqMap[topic]
			producer := producers[topic]
			_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
				Payload: util.CreatePing(seq).Serialize(msg_size),
			})

			if err != nil {
				panic("publish error:" + err.Error())
			}
			seq++
			seqMap[topic] = seq
			cnt++
		}
		fmt.Printf("%s: %d ns. %d times. %f ns/op( %f op/s )\n", "", et-st, cnt, float64(et-st)/float64(cnt), float64(cnt)*1e9/float64(et-st))
	}
	println("end-benchmark")

}
