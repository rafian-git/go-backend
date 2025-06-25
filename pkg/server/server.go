package server

import (
	"context"
	"gitlab.techetronventures.com/core/backend/pkg/trace_id"
	"net"
	"runtime/debug"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"

	grpcmid "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrec "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	ae "gitlab.techetronventures.com/core/backend/pkg/apierror"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"gitlab.techetronventures.com/core/backend/pkg/user_id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server initializes and wraps a gRPC server.
type Server struct {
	log  *log.Logger
	conf *Config
	l    net.Listener
	srv  *grpc.Server
	wg   sync.WaitGroup
}

// The Config represents gRPC server configurations.
type Config struct {
	Network              string `json:"network" yaml:"network" toml:"network" mapstructure:"network"`                                                             // nolint
	ServiceName          string `json:"service_name" yaml:"service_name" toml:"service_name" mapstructure:"service_name"`                                         // nolint
	Bind                 string `json:"bind" yaml:"bind" toml:"bind" mapstructure:"bind"`                                                                         // nolint
	MaxRecvMsgSize       int    `json:"max_recv_msg_size" yaml:"max_recv_msg_size" toml:"max_recv_msg_size" mapstructure:"max_recv_msg_size"`                     // nolint
	MaxSendMsgSize       int    `json:"max_send_msg_size" yaml:"max_send_msg_size" toml:"max_send_msg_size" mapstructure:"max_send_msg_size"`                     // nolint
	ReadBufferSize       int    `json:"read_buffer_size" yaml:"read_buffer_size" toml:"read_buffer_size" mapstructure:"read_buffer_size"`                         // nolint
	WriteBufferSize      int    `json:"write_buffer_size" yaml:"write_buffer_size" toml:"write_buffer_size" mapstructure:"write_buffer_size"`                     // nolint
	Development          bool   `json:"development" yaml:"development" toml:"development" mapstructure:"development"`                                             // nolint
	PayloadLogs          bool   `json:"payload_logs" yaml:"payload_logs" toml:"payload_logs" mapstructure:"payload_logs"`                                         // nolint
	MaxConcurrentStreams uint32 `json:"max_concurrent_streams" yaml:"max_concurrent_streams" toml:"max_concurrent_streams" mapstructure:"max_concurrent_streams"` // nolint
}

// NewConfig returns default configurations.
func NewConfig() (c *Config) {
	c = new(Config)
	c.Network = "tcp"
	// the address to bind is not defined by default
	c.MaxRecvMsgSize = 4194304
	c.MaxSendMsgSize = 4194304
	c.ReadBufferSize = 4194304
	c.WriteBufferSize = 4194304
	// log requests and responses payloads by default
	c.PayloadLogs = true
	c.MaxConcurrentStreams = 100
	return
}

func (s *Server) recoveryHandlerFunc(ctx context.Context, p interface{}) error {

	// s.log.Error(ctx, "panic", zap.Any("recover", p))
	debug.PrintStack()
	return ae.New(ae.Internal, "Something went wrong")
}

func (s *Server) unaryInterceptors(log *log.Logger, reg *prometheus.Registry, ch ...grpc.UnaryServerInterceptor) ( // nolint
	uis grpc.ServerOption) {

	var chain = []grpc.UnaryServerInterceptor{}

	chain = append(chain, ch...)

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	// Create a customized counter metric.
	customizedCounterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})

	histogramMetrics := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      s.conf.ServiceName,
		Help:      "histogram distribution",
		Namespace: "",
	})
	grpcMetrics.EnableHandlingTimeHistogram()

	reg.MustRegister(grpcMetrics, customizedCounterMetric, histogramMetrics)
	customizedCounterMetric.WithLabelValues("Test")

	// histogramMetrics.WithLabelValues("test_histogram")

	// if a != nil {
	// 	chain = append(chain, a.UnaryServerInterceptor())
	// }

	chain = append(chain,
		grpcrec.UnaryServerInterceptor(
			grpcrec.WithRecoveryHandlerContext(s.recoveryHandlerFunc),
		), user_id.UserIdUnaryInterceptor, trace_id.UnaryInterceptor, grpc.UnaryServerInterceptor(grpcMetrics.UnaryServerInterceptor()))

	if s.conf.PayloadLogs {
		if s.conf.PayloadLogs {
			chain = append(chain, grpczap.PayloadUnaryServerInterceptor(
				log.Named("payload").Logger,
				alwaysTrueDecider,
			))
		}
	}

	return grpc.UnaryInterceptor(grpcmid.ChainUnaryServer(chain...))
}

func (s *Server) options(log *log.Logger, reg *prometheus.Registry, ch ...grpc.UnaryServerInterceptor) (
	opts []grpc.ServerOption) {

	if s.conf.MaxRecvMsgSize > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(s.conf.MaxRecvMsgSize))
	}
	if s.conf.MaxSendMsgSize > 0 {
		opts = append(opts, grpc.MaxSendMsgSize(s.conf.MaxSendMsgSize))
	}
	if s.conf.ReadBufferSize > 0 {
		opts = append(opts, grpc.ReadBufferSize(s.conf.ReadBufferSize))
	}
	if s.conf.WriteBufferSize > 0 {
		opts = append(opts, grpc.WriteBufferSize(s.conf.WriteBufferSize))
	}

	if s.conf.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(s.conf.MaxConcurrentStreams))
	}

	if uis := s.unaryInterceptors(log, reg, ch...); uis != nil {
		opts = append(opts, uis)
	}
	// if sis := s.streamInterceptors(a, log); sis != nil {
	// 	opts = append(opts, sis)
	// }
	// TODO: add unary & stream interceptors, including
	//  - [ ] opentracing (optional, low priority)
	//  - [ ] metrics  (optional, low priority)
	// if s.conf.Tracing.Enable {
	// 	tracer, err := tracing.NewZipkinTracer(s.conf.Tracing.EndpointURL, s.conf.Tracing.ServiceName, s.conf.Tracing.EndpointURL)
	// 	if err != nil {
	// 		s.log.Error(context.Background(), err.Error())
	// 	}
	//
	// 	opentracing.SetGlobalTracer(zipkinot.Wrap(tracer))
	//
	// 	opts = append(opts, grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)))
	// }

	return
}

// New gRPC server instance with all required middleware and listener.
func New(log *log.Logger, conf *Config, registry *prometheus.Registry, ui ...grpc.UnaryServerInterceptor) (s *Server, err error) {
	s = new(Server)
	s.conf = conf
	s.log = log.Named("grpc")
	if s.l, err = net.Listen(conf.Network, conf.Bind); err != nil {
		return
	}
	s.srv = grpc.NewServer(s.options(s.log, registry, ui...)...)
	return
}

// Start the Server. It blocks. On grpc.ErrServerStopped it returns nil.
func (s *Server) Start() (err error) {
	reflection.Register(s.srv)
	if err = s.srv.Serve(s.l); err != nil && err != grpc.ErrServerStopped {
		panic(err)
		return
	}
	return nil // the ServerStopped error, not a real error
}

func (s *Server) run() {
	defer s.wg.Done()
	if err := s.srv.Serve(s.l); err != nil && err != grpc.ErrServerStopped {
		panic(err)
	}
}

// Run server in goroutine. Panic on error (excluding grpc.ErrServerStopped).
func (s *Server) Run() {
	reflection.Register(s.srv)
	s.wg.Add(1)
	go s.run()
}

// GRPC server instance. Use this method to register a service implementation.
func (s *Server) GRPC() *grpc.Server {
	return s.srv
}

// Close the Server gracefully.
func (s *Server) Close() (err error) {
	s.srv.GracefulStop()
	err = s.l.Close()
	s.wg.Wait()
	return
}

// for grpc logging
func alwaysTrueDecider(context.Context, string, interface{}) bool {
	return true
}
