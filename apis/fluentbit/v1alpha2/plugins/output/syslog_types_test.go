package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	. "github.com/onsi/gomega"
)

func TestOutput_Syslog_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test-namespace")

	// Initialize Syslog struct based on your provided types
	syslog := Syslog{
		Host:               "127.0.0.1",
		Port:               utils.ToPtr[int32](514),
		Mode:               "tcp",
		SyslogFormat:       "rfc5424",
		SyslogMaxSize:      utils.ToPtr[int32](2048),
		SyslogSeverityKey:  "severity",
		SyslogFacilityKey:  "facility",
		SyslogHostnameKey:  "hostname",
		SyslogAppnameKey:   "appname",
		SyslogProcessIDKey: "pid",
		SyslogMessageIDKey: "msgid",
		SyslogSDKey:        "structured_data",
		SyslogMessageKey:   "message",
		TotalLimitSize:     "1G",
		workers:            utils.ToPtr[int32](2),
	}

	// Define the expected Key-Value pairs based on your Params() logic
	expected := params.NewKVs()
	expected.Insert("Host", "127.0.0.1")
	expected.Insert("port", "514")
	expected.Insert("mode", "tcp")
	expected.Insert("syslog_hostname_key", "hostname")
	expected.Insert("syslog_appname_key", "appname")
	expected.Insert("syslog_message_key", "message")
	expected.Insert("syslog_format", "rfc5424")
	expected.Insert("syslog_severity_key", "severity")
	expected.Insert("syslog_facility_key", "facility")
	expected.Insert("syslog_procid_key", "pid") // Note: matches your implementation
	expected.Insert("syslog_msgid_key", "msgid")
	expected.Insert("syslog_sd_key", "structured_data") // Note: matches your implementation
	expected.Insert("storage.total_limit_size", "1G")
	expected.Insert("syslog_maxsize", "2048")
	expected.Insert("workers", "2")

	// Execute the translation
	kvs, err := syslog.Params(sl)

	// Assertions
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
