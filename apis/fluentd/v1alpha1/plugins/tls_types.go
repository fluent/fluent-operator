package plugins

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Fluentd provides integrated support for Transport Layer Security (TLS) and it predecessor Secure Sockets Layer (SSL) respectively.
type TLS struct {
	// Disable certificate validation
	Insecure *bool `json:"insecure,omitempty"`
	// Absolute path to CA certificate file
	CAFile string `json:"caFile,omitempty"`
	// Absolute path to Certificate file
	CRTFile string `json:"crtFile,omitempty"`
	// Absolute path to private Key file
	KeyFile string `json:"keyFile,omitempty"`
}

type TLSLoader struct {
	client    client.Client
	namespace string
}

func NewTLSMapLoader(c client.Client, ns string) TLSLoader {
	return TLSLoader{
		client:    c,
		namespace: ns,
	}
}
