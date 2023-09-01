package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_Gelf_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	dd := Gelf{
		Host:            "127.0.0.1",
		Port:            ptrInt32(1234),
		Mode:            "udp",
		ShortMessageKey: "short_message",
		TimestampKey:    "timestamp",
		HostKey:         "host",
		FullMessageKey:  "full_message",
		LevelKey:        "level",
		PacketSize:      ptrInt32(1000),
		Compress:        ptrBool(true),
	}

	expected := params.NewKVs()
	expected.Insert("Host", "127.0.0.1")
	expected.Insert("Port", "1234")
	expected.Insert("Mode", "udp")
	expected.Insert("Gelf_Short_Message_Key", "short_message")
	expected.Insert("Gelf_Timestamp_Key", "timestamp")
	expected.Insert("Gelf_Host_Key", "host")
	expected.Insert("Gelf_Full_Message_Key", "full_message")
	expected.Insert("Gelf_Level_Key", "level")
	expected.Insert("Packet_Size", "1000")
	expected.Insert("Compress", "true")

	kvs, err := dd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func ptrInt32(v int32) *int32 {
	return &v
}
