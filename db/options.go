package db

import (
	"time"
)

type options struct {
	// 本地队列长度
	jobChanLen int
	// 消费者数量
	consumerNum int
	// 生产者数量
	producerNum int
}

func (o *options) fix() {
	if o.jobChanLen <= 0 {
		o.jobChanLen = 1
	}
	if o.consumerNum <= 0 {
		o.consumerNum = 1
	}
	if o.producerNum <= 0 {
		o.producerNum = 1
	}
}

type Option func(*options)

func WithJobChanLen(l int) Option {
	return func(o *options) {
		o.jobChanLen = l
	}
}

func WithConsumerNum(num int) Option {
	return func(o *options) {
		o.consumerNum = num
	}
}
func WithBackoff(initialDelay, maxDelay time.Duration, maxRetries int, factor float64) Option {
	return func(o *options) {
	}
}
