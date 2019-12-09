package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/siddontang/go-log/log"
)

type sender struct {
	logWriter *kafka.Writer
	wrapMsg   func(*kafka.Message, interface{})
}

type recver struct {
	logWriter *kafka.Reader
}

type Sender interface {
	Send(...interface{}) error
	SetWrapper(fn func(*kafka.Message, interface{}))
}

type Recver interface {
	// Send(...interface{}) error
}

func NewSender(topic string, brokers []string) Sender {
	s := &sender{}

	s.logWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 5 * time.Millisecond,
	})
	return s
}

func (s *sender) SetWrapper(fn func(*kafka.Message, interface{})) {
	s.wrapMsg = fn
}

func (s *sender) Send(msgs ...interface{}) error {
	var list = make([]kafka.Message, len(msgs))
	for i := range msgs {
		data, err := json.Marshal(msgs[i])
		if err != nil {
			log.Println("marshal: ", err)
			return err
		}
		list[i].Value = data
		if s.wrapMsg != nil {
			s.wrapMsg(&list[i], msgs[i])
		}
	}
	return s.logWriter.WriteMessages(context.Background(), list...)
}

func NewRecver(topic string, brokers []string) Recver {
	s := &recver{}

	s.logWriter = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return s
}
