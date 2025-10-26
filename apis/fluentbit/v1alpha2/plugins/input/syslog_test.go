package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	. "github.com/onsi/gomega"
)

func TestSyslog_Name(t *testing.T) {
	g := NewGomegaWithT(t)
	syslog := Syslog{}
	g.Expect(syslog.Name()).To(Equal("syslog"))
}

func TestSyslog_Params(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "test namespace")

	syslog := Syslog{
		Mode:              "tcp",
		Listen:            "0.0.0.0",
		Port:              utils.ToPtr[int32](514),
		Path:              "/tmp/syslog.sock",
		UnixPerm:          utils.ToPtr[int32](644),
		Parser:            "syslog-rfc5424",
		BufferChunkSize:   "32KB",
		BufferMaxSize:     "256KB",
		ReceiveBufferSize: "1MB",
		SourceAddressKey:  "source_address",
		Tag:               "syslog.tag",
	}

	expected := params.NewKVs()
	expected.Insert("Mode", "tcp")
	expected.Insert("Listen", "0.0.0.0")
	expected.Insert("Path", "/tmp/syslog.sock")
	expected.Insert("Parser", "syslog-rfc5424")
	expected.Insert("Buffer_Chunk_Size", "32KB")
	expected.Insert("Buffer_Max_Size", "256KB")
	expected.Insert("Receive_Buffer_Size", "1MB")
	expected.Insert("Source_Address_Key", "source_address")
	expected.Insert("Tag", "syslog.tag")
	expected.Insert("Port", "514")
	expected.Insert("Unix_Perm", "644")

	kvs, err := syslog.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
