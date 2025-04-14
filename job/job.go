package job

import "sync"

type JobServer interface {
	ConsumerNum() int
	ProducerNum() int
	Consume()
	Product()
	Close()
}

type JobServerImpl struct {
	WG     *sync.WaitGroup
	Resume bool
}

func DefaultJobServerImpl() *JobServerImpl {
	return &JobServerImpl{
		WG:     &sync.WaitGroup{},
		Resume: true,
	}
}
