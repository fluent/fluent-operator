package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHTTP_Params(t *testing.T) {
	g := NewGomegaWithT(t)
	fcb := fake.ClientBuilder{}
	fc := fcb.WithObjects(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test_namespace", Name: "http_secret"},
		Data: map[string][]byte{
			"http_user":   []byte("expected_http_user"),
			"http_passwd": []byte("expected_http_passwd"),
		},
	}).Build()

	sl := plugins.NewSecretLoader(fc, "test_namespace")
	h := HTTP{
		Host:           "example.com",
		Port:           utils.ToPtr[int32](443),
		HTTPUser:       &plugins.Secret{ValueFrom: plugins.ValueSource{SecretKeyRef: v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "http_secret"}, Key: "http_user"}}},
		HTTPPasswd:     &plugins.Secret{ValueFrom: plugins.ValueSource{SecretKeyRef: v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "http_secret"}, Key: "http_passwd"}}},
		Uri:            "/logs",
		Format:         "json",
		Headers:        map[string]string{"X-Custom": "value"},
		TLS:            &plugins.TLS{Verify: utils.ToPtr(false)},
		Networking:     &plugins.Networking{SourceAddress: utils.ToPtr("expected_source_address")},
		TotalLimitSize: "512M",
	}

	expected := params.NewKVs()
	expected.Insert("http_User", "expected_http_user")
	expected.Insert("http_Passwd", "expected_http_passwd")
	expected.Insert("host", "example.com")
	expected.Insert("port", "443")
	expected.Insert("uri", "/logs")
	expected.Insert("format", "json")
	expected.Insert("header", " X-Custom    value")
	expected.Insert("tls", "On")
	expected.Insert("tls.verify", "false")
	expected.Insert("net.source_address", "expected_source_address")
	expected.Insert("storage.total_limit_size", "512M")

	kvs, err := h.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
