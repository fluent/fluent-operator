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

func TestSplunk_Params(t *testing.T) {
	g := NewGomegaWithT(t)
	fcb := fake.ClientBuilder{}
	fc := fcb.WithObjects(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test_namespace", Name: "splunk_secret"},
		Data: map[string][]byte{
			"splunk_token": []byte("expected_splunk_token"),
		},
	}).Build()

	sl := plugins.NewSecretLoader(fc, "test_namespace")
	s := Splunk{
		Host:           "splunk.example.com",
		Port:           utils.ToPtr[int32](8088),
		SplunkToken:    &plugins.Secret{ValueFrom: plugins.ValueSource{SecretKeyRef: v1.SecretKeySelector{LocalObjectReference: v1.LocalObjectReference{Name: "splunk_secret"}, Key: "splunk_token"}}},
		TotalLimitSize: "512M",
	}

	expected := params.NewKVs()
	expected.Insert("splunk_token", "expected_splunk_token")
	expected.Insert("host", "splunk.example.com")
	expected.Insert("port", "8088")
	expected.Insert("storage.total_limit_size", "512M")

	kvs, err := s.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
