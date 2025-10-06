package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestOpenTelemetry_Params(t *testing.T) {
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
	ot := OpenTelemetry{
		Host:                  "otlp-collector.example.com",
		Port:                  ptrAny(int32(443)),
		HTTPUser:              &plugins.Secret{ValueFrom: plugins.ValueSource{SecretKeyRef: v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "http_secret"}, Key: "http_user"}}},
		HTTPPasswd:            &plugins.Secret{ValueFrom: plugins.ValueSource{SecretKeyRef: v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "http_secret"}, Key: "http_passwd"}}},
		Proxy:                 "expected_proxy",
		MetricsUri:            "expected_metrics_uri",
		LogsUri:               "expected_logs_uri",
		TracesUri:             "expected_traces_uri",
		Header:                map[string]string{"custom_header_key": "custom_header_val"},
		LogResponsePayload:    ptrBool(true),
		AddLabel:              map[string]string{"add_label_key": "add_label_val"},
		LogsBodyKeyAttributes: ptrBool(true),
		LogsBodyKey:           "expected_logs_body_key",
		TLS:                   &plugins.TLS{Verify: ptrBool(false)},
		Networking:            &plugins.Networking{SourceAddress: ptrAny("expected_source_address")},
	}

	expected := params.NewKVs()
	expected.Insert("host", "otlp-collector.example.com")
	expected.Insert("port", "443")
	expected.Insert("http_user", "expected_http_user")
	expected.Insert("http_passwd", "expected_http_passwd")
	expected.Insert("proxy", "expected_proxy")
	expected.Insert("metrics_uri", "expected_metrics_uri")
	expected.Insert("logs_uri", "expected_logs_uri")
	expected.Insert("traces_uri", "expected_traces_uri")
	expected.Insert(header, " custom_header_key    custom_header_val")
	expected.Insert("log_response_payload", "true")
	expected.Insert(addLabel, " add_label_key    add_label_val")
	expected.Insert("logs_body_key_attributes", "true")
	expected.Insert("logs_body_key", "expected_logs_body_key")
	expected.Insert("tls", "On")
	expected.Insert("tls.verify", "false")
	expected.Insert("net.source_address", "expected_source_address")

	kvs, err := ot.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
