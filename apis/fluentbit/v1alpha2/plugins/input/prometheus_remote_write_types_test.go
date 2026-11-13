package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInput_PrometheusRemoteWrite_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	prw := PrometheusRemoteWrite{
		Tag:                    "prom.rw",
		Listen:                 "0.0.0.0",
		Port:                   new(int32(9090)),
		URI:                    "/api/prom/push",
		BufferMaxSize:          "8M",
		BufferChunkSize:        "1M",
		SuccessfulResponseCode: new(int32(204)),
		Threaded:               new(true),
	}

	expected := params.NewKVs()
	expected.Insert("tag", "prom.rw")
	expected.Insert("listen", "0.0.0.0")
	expected.Insert("port", "9090")
	expected.Insert("uri", "/api/prom/push")
	expected.Insert("buffer_max_size", "8M")
	expected.Insert("buffer_chunk_size", "1M")
	expected.Insert("successful_response_code", "204")
	expected.Insert("threaded", "true")

	kvs, err := prw.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestInput_PrometheusRemoteWrite_ParamsOnlySetFields(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kvs, err := (&PrometheusRemoteWrite{}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(params.NewKVs()))

	kvs, err = (&PrometheusRemoteWrite{Port: new(int32(8080)), Threaded: new(false)}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	expected := params.NewKVs()
	expected.Insert("port", "8080")
	expected.Insert("threaded", "false")
	g.Expect(kvs).To(Equal(expected))
}

func TestInput_PrometheusRemoteWrite_ParamsTLS(t *testing.T) {
	g := NewGomegaWithT(t)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "remote-write-tls", Namespace: "test namespace"},
		Data: map[string][]byte{
			"key-password": []byte("s3cr3t"),
		},
	}
	sl := plugins.NewSecretLoader(fake.NewClientBuilder().WithObjects(secret).Build(), "test namespace")

	keyPassword := func(key string) *plugins.Secret {
		return &plugins.Secret{ValueFrom: plugins.ValueSource{
			SecretKeyRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: "remote-write-tls"},
				Key:                  key,
			},
		}}
	}

	prw := PrometheusRemoteWrite{
		Port: new(int32(8443)),
		TLS: &plugins.TLS{
			CRTFile:     "/fluent-bit/tls/tls.crt",
			KeyFile:     "/fluent-bit/tls/tls.key",
			KeyPassword: keyPassword("key-password"),
		},
	}

	expected := params.NewKVs()
	expected.Insert("port", "8443")
	expected.Insert("tls", "On")
	expected.Insert("tls.crt_file", "/fluent-bit/tls/tls.crt")
	expected.Insert("tls.key_file", "/fluent-bit/tls/tls.key")
	expected.Insert("tls.key_passwd", "s3cr3t")

	kvs, err := prw.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))

	prw.KeyPassword = keyPassword("missing")
	kvs, err = prw.Params(sl)
	g.Expect(err).To(HaveOccurred())
	g.Expect(kvs).To(BeNil())
}
