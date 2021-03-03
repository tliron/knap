module github.com/tliron/knap

go 1.16

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

require (
	github.com/gofrs/flock v0.8.0
	github.com/golang/protobuf v1.4.3
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/tebeka/atexit v0.3.0
	github.com/tliron/kutil v0.1.21
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)
