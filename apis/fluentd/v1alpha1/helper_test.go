package v1alpha1

import (
	"os"
	"testing"

	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/common"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/filter"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/input"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/output"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	forwardId       = "forward-001"
	forwardLogLevel = "info"
	forwardLabel    = "forward-test"

	transportTls          = "tls"
	forwardCertPath       = "/etc/td-agent/certs/fluentd.crt"
	forwardPrivateKeyPath = "/etc/td-agent/certs/fluentd.key"

	forwardPort int32 = 24224

	inputs = []input.Input{
		{
			InputCommon: input.InputCommon{
				Id:       &forwardId,
				LogLevel: &forwardLogLevel,
				Label:    &forwardLabel,
			},
			Forward: &input.Forward{
				Transport: &common.Transport{
					Protocol:       &transportTls,
					CertPath:       &forwardCertPath,
					PrivateKeyPath: &forwardPrivateKeyPath,
				},
				Port: &forwardPort,
			},
		},
	}

	recordKey1   = "avg"
	recordValue1 = `${record["total"] / record["count"]}`
	recordKey2   = "message"
	recordValue2 = `yay, ${record["message"]}`

	regexpKey1      = "message"
	regexpPattern1  = "/cool/"
	regexpKey2      = "hostname"
	regexpPattern2  = `/^web\d+\.example\.com$/`
	excludeKey1     = "message"
	excludePattern1 = "/exclude/"
	excludeKey2     = "hostname"
	excludePattern2 = `/^web\d+\.error\.com$/`

	regexpParser     = "regexp"
	regexpExpression = `/^(?<host>[^ ]*) [^ ]* (?<user>[^ ]*) [(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^ ]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)$/`

	timeFormat = `%d/%b/%Y:%H:%M:%S %z`
	formatType = "json"

	records = []*filter.Record{
		{
			Key:   &recordKey1,
			Value: &recordValue1,
		},
		{
			Key:   &recordKey2,
			Value: &recordValue2,
		},
	}

	bufferId             = "common_buffer"
	bufferType           = "file"
	bufferPath           = "/fluentd/buffer/"
	bufferTag            = "buffertag.*"
	bufferTotalLimitSize = "5GB"

	buffer = common.Buffer{
		BufferCommon: common.BufferCommon{
			Id:   &bufferId,
			Type: &bufferType,
		},
		Path:           &bufferPath,
		Tag:            bufferTag,
		TotalLimitSize: &bufferTotalLimitSize,
	}

	serverHost = "host"
	serverPort = "14423"
	sdType     = "file"
	sdPath     = "/sd/path"

	servers = []*common.Server{
		{
			Host: &serverHost,
			Port: &serverPort,
		},
	}

	serversDiscovery = common.ServiceDiscovery{
		SDCommon: common.SDCommon{Type: &sdType},
		FileServiceDiscovery: &common.FileServiceDiscovery{
			Path: &sdPath,
		},
	}

	authMethod = "basic"
	username   = "username"
	password   = "password"

	auth = common.Auth{
		Method:   &authMethod,
		Username: &username,
		Password: &password,
	}

	endpoint           = "http://logserver.com:9000/api"
	opentimeout uint16 = 2

	brokers = "broker1,broker2"
)

func Test_RenderMainConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace", nil)

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label2": "lv2",
		"label3": "lv3",
	}

	filterspec1, filterspec2, filterspec3 := createFilterSpecs()
	outputspec1, outputspec2, outputspec3 := createOutputSpecs()

	clusterFilters := []ClusterFilter{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter1",
			},
			Spec: ClusterFilterSpec(filterspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter2",
			},
			Spec: ClusterFilterSpec(filterspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter3",
			},
			Spec: ClusterFilterSpec(filterspec3),
		},
	}

	clusterOutputs := []ClusterOutput{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput1",
			},
			Spec: ClusterOutputSpec(outputspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput2",
			},
			Spec: ClusterOutputSpec(outputspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput3",
			},
			Spec: ClusterOutputSpec(outputspec3),
		},
	}

	filters := []Filter{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter1",
				Namespace: "testnamespace",
			},
			Spec: FilterSpec(filterspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter2",
				Namespace: "testnamespace",
			},
			Spec: FilterSpec(filterspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter3",
				Namespace: "testnamespace",
			},
			Spec: FilterSpec(filterspec3),
		},
	}

	outputs := []Output{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output1",
				Namespace: "testnamespace",
			},
			Spec: OutputSpec(outputspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output2",
				Namespace: "testnamespace",
			},
			Spec: OutputSpec(outputspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output3",
				Namespace: "testnamespace",
			},
			Spec: OutputSpec(outputspec3),
		},
	}

	clustercfg := ClusterFluentdConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentd.fluent.io/v1alpha1",
			Kind:       "ClusterFluentdConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "clusterfluentdconfig0",
		},
		Spec: ClusterFluentdConfigSpec{
			WatchedNamespaces: []string{"ns1", "ns2"},
			WatchedLabels:     labels,
		},
	}

	cfg := FluentdConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentd.fluent.io/v1alpha1",
			Kind:       "FluentdConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentdconfig0",
			Namespace: "testnamespace",
		},
		Spec: FluentdConfigSpec{
			WatchedLabels: labels,
		},
	}

	psr := NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, inputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&clustercfg)
	g.Expect(err).NotTo(HaveOccurred())

	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, clustercfg.GetCfgId(), clusterFilters, clusterOutputs)
	psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)

	cfgRouter, err := psr.BuildCfgRouter(&cfg)
	g.Expect(err).NotTo(HaveOccurred())

	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, cfg.GetCfgId(), filters, outputs)
	psr.WithCfgResources(*cfgRouter.Label, cfgResources)

	expected, _ := os.ReadFile("./tests/expected-main-cfgs.cfg")

	// we should not see any permutations in serialized config
	i := 0
	for i < 5 {
		config, errs := psr.RenderMainConfig()
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(expected)).To(Equal(config))

		i++
	}

}

func createFilterSpecs() (FilterSpec, FilterSpec, FilterSpec) {
	filterSpec1 := FilterSpec{
		Filters: []filter.Filter{
			{
				RecordTransformer: &filter.RecordTransformer{
					Records: records,
				},
			},
			{
				Grep: &filter.Grep{
					Regexps: []*filter.Regexp{
						{
							Key:     &regexpKey1,
							Pattern: &regexpPattern1,
						},
						{
							Key:     &regexpKey2,
							Pattern: &regexpPattern2,
						},
					},
					Excludes: []*filter.Exclude{
						{
							Key:     &excludeKey1,
							Pattern: &excludePattern1,
						},
						{
							Key:     &excludeKey2,
							Pattern: &excludePattern2,
						},
					},
				},
			},
		},
	}

	filterSpec2 := FilterSpec{
		Filters: []filter.Filter{
			{
				RecordTransformer: &filter.RecordTransformer{
					Records: records,
				},
			},
			{
				Stdout: &filter.Stdout{
					Format: &common.Format{
						FormatCommon: common.FormatCommon{
							Type: &formatType,
						},
					},
				},
			},
		},
	}

	filterSpec3 := FilterSpec{
		Filters: []filter.Filter{
			{
				Parser: &filter.Parser{
					Parse: &common.Parse{
						ParseCommon: common.ParseCommon{
							Type: &regexpParser,
						},
						Time: common.Time{
							TimeFormat: &timeFormat,
						},
						Expression: &regexpExpression,
					},
				},
			},
			{
				Stdout: &filter.Stdout{
					Format: &common.Format{
						FormatCommon: common.FormatCommon{
							Type: &formatType,
						},
					},
				},
			},
		},
	}

	return filterSpec1, filterSpec2, filterSpec3
}

func createOutputSpecs() (OutputSpec, OutputSpec, OutputSpec) {
	outputSpec1 := OutputSpec{
		Outputs: []output.Output{
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Forward: &output.Forward{
					Servers:          servers,
					ServiceDiscovery: &serversDiscovery,
				},
			},
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Kafka: &output.Kafka2{
					Brokers: &brokers,
				},
			},
		},
	}

	outputSpec2 := OutputSpec{
		Outputs: []output.Output{
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Forward: &output.Forward{
					Servers:          servers,
					ServiceDiscovery: &serversDiscovery,
				},
			},
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Http: &output.Http{
					Auth:        &auth,
					Endpoint:    &endpoint,
					OpenTimeout: &opentimeout,
				},
			},
		},
	}

	outputSpec3 := OutputSpec{
		Outputs: []output.Output{
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Http: &output.Http{
					Auth:        &auth,
					Endpoint:    &endpoint,
					OpenTimeout: &opentimeout,
				},
			},
			{
				Match: common.Match{
					Buffer: &buffer,
				},
				Kafka: &output.Kafka2{
					Brokers: &brokers,
				},
			},
		},
	}

	return outputSpec1, outputSpec2, outputSpec3
}
