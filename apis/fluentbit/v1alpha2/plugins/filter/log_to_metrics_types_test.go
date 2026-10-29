package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestFilter_LogToMetrics_FlushInterval(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	l := LogToMetrics{
		Tag:               "metrics",
		MetricMode:        "counter",
		MetricName:        "count_all",
		FlushIntervalSec:  new(int32(30)),
		FlushIntervalNsec: new(int32(500000000)),
	}

	expected := params.NewKVs()
	expected.Insert("Tag", "metrics")
	expected.Insert("Metric_mode", "counter")
	expected.Insert("Metric_name", "count_all")
	expected.Insert("Flush_Interval_Sec", "30")
	expected.Insert("Flush_Interval_Nsec", "500000000")

	kvs, err := l.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestFilter_LogToMetrics_NoFlushInterval(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	l := LogToMetrics{Tag: "metrics"}

	expected := params.NewKVs()
	expected.Insert("Tag", "metrics")

	kvs, err := l.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
