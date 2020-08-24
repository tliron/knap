module github.com/tliron/knap

go 1.15

replace github.com/tliron/puccini => /Depot/Projects/RedHat/puccini

replace github.com/tliron/turandot => /Depot/Projects/RedHat/turandot

require (
	github.com/google/uuid v1.1.1
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20200626054723-37f83d1996bc
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/spf13/cobra v1.0.0
	github.com/tebeka/atexit v0.3.0
	github.com/tliron/puccini v0.0.0-00010101000000-000000000000
	github.com/tliron/turandot v0.0.0-00010101000000-000000000000
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/code-generator v0.18.8 // indirect
)
