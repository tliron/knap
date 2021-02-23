package provider

import (
	"context"
	"net"
	"time"

	"github.com/tliron/kutil/logging"
	"google.golang.org/grpc"
)

//
// Client
//

type Client struct {
	socketName string
	connection *grpc.ClientConn
	client     ProviderClient
	log        logging.Logger
}

func NewClient(socketName string, log logging.Logger) (*Client, error) {
	log.Infof("creating provider gRPC client on socket %s", socketName)
	if connection, err := grpc.Dial(socketName, grpc.WithInsecure(), grpc.WithDialer(dialer)); err == nil {
		return &Client{
			socketName: socketName,
			connection: connection,
			client:     NewProviderClient(connection),
			log:        log,
		}, nil
	} else {
		return nil, err
	}
}

func (self *Client) Release() {
	self.log.Infof("closing provider gRPC client connection on socket %s", self.socketName)
	self.connection.Close()
}

func (self *Client) CreateCniConfig(name string, hints map[string]string) (string, error) {
	self.log.Infof("calling gRPC CreateCniConfig on socket %s", self.socketName)

	context, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := CreateCniConfigRequest{
		Name:  name,
		Hints: hints,
	}

	if reply, err := self.client.CreateCniConfig(context, &request); err == nil {
		return reply.Config, nil
	} else {
		return "", err
	}
}

func dialer(address string, duration time.Duration) (net.Conn, error) {
	if unixAddress, err := net.ResolveUnixAddr("unix", address); err == nil {
		return net.DialUnix("unix", nil, unixAddress)
	} else {
		return nil, err
	}
}
