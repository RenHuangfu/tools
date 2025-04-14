package job

import (
	"context"
	"net/url"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	name     string
	endpoint string

	srvs []JobServer

	debug   bool
	started bool
	ctx     context.Context
	wg      *sync.WaitGroup
}

// Deprecated: pubsub instead
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		ctx:     context.Background(),
		started: false,
		wg:      &sync.WaitGroup{},
	}
	for _, o := range opts {
		o(srv)
	}

	return srv
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Start(ctx context.Context) error {
	if s.started {
		return nil
	}

	for i := range s.srvs {
		for j := 0; j < s.srvs[i].ConsumerNum(); j++ {
			s.wg.Add(1)
			go func(srv JobServer) {
				defer s.wg.Done()
				srv.Consume()
			}(s.srvs[i])
		}

		for k := 0; k < s.srvs[i].ProducerNum(); k++ {
			s.wg.Add(1)
			go func(srv JobServer) {
				defer s.wg.Done()
				srv.Product()
			}(s.srvs[i])
		}
	}

	log.Infof("[%s] server listening on: %s", s.name, s.endpoint)

	s.ctx = ctx
	s.started = true

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	log.Infof("[%s] server stopping", s.name)
	for i := range s.srvs {
		go s.srvs[i].Close()
	}
	s.wg.Wait()
	s.started = false
	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	return &url.URL{Host: s.endpoint}, nil
}

func (s *Server) Register(srvs ...JobServer) {
	s.srvs = append(s.srvs, srvs...)
}
