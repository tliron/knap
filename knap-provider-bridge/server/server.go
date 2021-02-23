package server

import (
	"context"
	"net"

	"github.com/tliron/knap/provider"
	"github.com/tliron/kutil/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//
// Server
//

type Server struct {
	StateFilename string
	SocketName    string
	Log           logging.Logger
}

func NewServer(stateFilename string, socketName string, log logging.Logger) *Server {
	return &Server{
		StateFilename: stateFilename,
		SocketName:    socketName,
		Log:           log,
	}
}

func (self *Server) Start() error {
	if err := self.InitializeState(); err == nil {
		self.Log.Infof("starting provider gRPC server on socket %s", self.SocketName)
		if listener, err := net.Listen("unix", self.SocketName); err == nil {
			server := grpc.NewServer()
			provider.RegisterProviderServer(server, self)
			reflection.Register(server)
			return server.Serve(listener)
		} else {
			return err
		}
	} else {
		return err
	}
}

// provider.ProviderServer interface
func (self *Server) CreateCniConfig(context context.Context, request *provider.CreateCniConfigRequest) (*provider.CreateCniConfigReply, error) {
	self.Log.Infof("gRPC CreateCniConfig called on socket %s", self.SocketName)
	if config, err := self.CreateBridgeCniConfig(request.Name); err == nil {
		return &provider.CreateCniConfigReply{
			Config: config,
		}, nil
	} else {
		return nil, err
	}
}
