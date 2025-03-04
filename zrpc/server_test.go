package zrpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/l306287405/go-zero/core/discov"
	"github.com/l306287405/go-zero/core/logx"
	"github.com/l306287405/go-zero/core/service"
	"github.com/l306287405/go-zero/core/stat"
	"github.com/l306287405/go-zero/core/stores/redis"
	"github.com/l306287405/go-zero/zrpc/internal"
	"github.com/l306287405/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
)

func TestServer_setupInterceptors(t *testing.T) {
	server := new(mockedServer)
	err := setupInterceptors(server, RpcServerConf{
		Auth: true,
		Redis: redis.RedisKeyConf{
			RedisConf: redis.RedisConf{
				Host: "any",
				Type: redis.NodeType,
			},
			Key: "foo",
		},
		CpuThreshold: 10,
		Timeout:      100,
	}, new(stat.Metrics))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(server.unaryInterceptors))
	assert.Equal(t, 1, len(server.streamInterceptors))
}

func TestServer(t *testing.T) {
	SetServerSlowThreshold(time.Second)
	srv := MustNewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn:      "localhost:8080",
		Etcd:          discov.EtcdConf{},
		Auth:          false,
		Redis:         redis.RedisKeyConf{},
		StrictControl: false,
		Timeout:       0,
		CpuThreshold:  0,
	}, func(server *grpc.Server) {
	})
	srv.AddOptions(grpc.ConnectionTimeout(time.Hour))
	srv.AddUnaryInterceptors(serverinterceptors.UnaryCrashInterceptor)
	srv.AddStreamInterceptors(serverinterceptors.StreamCrashInterceptor)
	go srv.Start()
	srv.Stop()
}

func TestServerError(t *testing.T) {
	_, err := NewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn: "localhost:8080",
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost"},
		},
		Auth:  true,
		Redis: redis.RedisKeyConf{},
	}, func(server *grpc.Server) {
	})
	assert.NotNil(t, err)
}

func TestServer_HasEtcd(t *testing.T) {
	srv := MustNewServer(RpcServerConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				ServiceName: "foo",
				Mode:        "console",
			},
		},
		ListenOn: "localhost:8080",
		Etcd: discov.EtcdConf{
			Hosts: []string{"notexist"},
			Key:   "any",
		},
		Redis: redis.RedisKeyConf{},
	}, func(server *grpc.Server) {
	})
	srv.AddOptions(grpc.ConnectionTimeout(time.Hour))
	srv.AddUnaryInterceptors(serverinterceptors.UnaryCrashInterceptor)
	srv.AddStreamInterceptors(serverinterceptors.StreamCrashInterceptor)
	go srv.Start()
	srv.Stop()
}

type mockedServer struct {
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

func (m *mockedServer) AddOptions(options ...grpc.ServerOption) {
}

func (m *mockedServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	m.streamInterceptors = append(m.streamInterceptors, interceptors...)
}

func (m *mockedServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	m.unaryInterceptors = append(m.unaryInterceptors, interceptors...)
}

func (m *mockedServer) SetName(s string) {
}

func (m *mockedServer) Start(register internal.RegisterFn) error {
	return nil
}
