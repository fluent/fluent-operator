package output

import (
	"strings"
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/common"
)

func TestForwardTlsCertLogicalStoreName(t *testing.T) {
	id := "forward-windows"
	host := "fluentd-aggregator"
	store := "Trust"
	o := &Output{
		Id: &id,
		Forward: &Forward{
			Servers:                 []*common.Server{{Host: &host}},
			TlsCertLogicalStoreName: &store,
		},
	}

	ps, err := o.Params(nil)
	if err != nil {
		t.Fatalf("Params returned an error: %v", err)
	}

	if got := ps.Store["tls_cert_logical_store_name"]; got != store {
		t.Errorf("tls_cert_logical_store_name = %q, want %q in %q", got, store, ps.String())
	}
	if strings.Contains(ps.String(), "tls_cert_logical_storeName") {
		t.Errorf("rendered config uses a parameter out_forward does not know: %q", ps.String())
	}
}
