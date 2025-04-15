package db

import (
	"context"
	"fmt"
	"github.com/RenHuangfu/tools/job"
)

type Broker[MSG, DB any] struct {
	job.JobServerImpl
	options

	db DB

	jobChan      chan MSG
	producerFunc func(ctx context.Context, db DB) ([]MSG, bool)
	consumerFunc func(ctx context.Context, msg MSG) error
}

// 使用producerFunc获取任务，打到内置的队列中，再由Consumer发送到srv
func New[MSG, DB any](
	db DB,
	producerFunc func(context.Context, DB) ([]MSG, bool),
	consumerFunc func(context.Context, MSG) error,
	opts ...Option,
) *Broker[MSG, DB] {
	o := &options{}
	o.fix()
	return &Broker[MSG, DB]{
		JobServerImpl: *job.DefaultJobServerImpl(),
		options:       *o,
		db:            db,
		jobChan:       make(chan MSG, o.jobChanLen),
		producerFunc:  producerFunc,
		consumerFunc:  consumerFunc,
	}
}

func (j *Broker[MSG, DB]) ConsumerNum() int {
	return j.consumerNum
}

func (j *Broker[MSG, DB]) ProducerNum() int {
	return j.producerNum
}

func (j *Broker[MSG, DB]) Consume() {
	j.WG.Add(1)
	defer j.WG.Done()

	for job := range j.jobChan {
		ctx := context.Background()
		err := j.consumerFunc(ctx, job)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (j *Broker[MSG, DB]) Product() {
	j.WG.Add(1)
	defer j.WG.Done()
	defer close(j.jobChan)

	for j.Resume {
		var msgs []MSG
		for _, msg := range msgs {
			j.jobChan <- msg
		}
	}
}

func (j *Broker[MSG, DB]) Close() {
	j.Resume = false
	j.WG.Wait()
}
