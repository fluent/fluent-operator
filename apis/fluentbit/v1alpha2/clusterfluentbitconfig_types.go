/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluentBitConfigSpec defines the desired state of ClusterFluentBitConfig
type FluentBitConfigSpec struct {
	// Service defines the global behaviour of the Fluent Bit engine.
	Service *Service `json:"service,omitempty"`
	// Select input plugins
	InputSelector metav1.LabelSelector `json:"inputSelector,omitempty"`
	// Select filter plugins
	FilterSelector metav1.LabelSelector `json:"filterSelector,omitempty"`
	// Select output plugins
	OutputSelector metav1.LabelSelector `json:"outputSelector,omitempty"`
	// Select parser plugins
	ParserSelector metav1.LabelSelector `json:"parserSelector,omitempty"`
	// Select multiline parser plugins
	MultilineParserSelector metav1.LabelSelector `json:"multilineParserSelector,omitempty"`
	// If namespace is defined, then the configmap and secret for fluent-bit is in this namespace.
	// If it is not defined, it is in the namespace of the fluentd-operator
	Namespace *string `json:"namespace,omitempty"`
}

type Storage struct {
	// Select an optional location in the file system to store streams and chunks of data/
	Path string `json:"path,omitempty"`
	// Configure the synchronization mode used to store the data into the file system
	// +kubebuilder:validation:Enum:=normal;full
	Sync string `json:"sync,omitempty"`
	// Enable the data integrity check when writing and reading data from the filesystem
	// +kubebuilder:validation:Enum:=on;off
	Checksum string `json:"checksum,omitempty"`
	// This option configure a hint of maximum value of memory to use when processing these records
	BacklogMemLimit string `json:"backlogMemLimit,omitempty"`
	// If the input plugin has enabled filesystem storage type, this property sets the maximum number of Chunks that can be up in memory
	MaxChunksUp *int64 `json:"maxChunksUp,omitempty"`
	// If http_server option has been enabled in the Service section, this option registers a new endpoint where internal metrics of the storage layer can be consumed
	// +kubebuilder:validation:Enum:=on;off
	Metrics string `json:"metrics,omitempty"`
	// When enabled, irrecoverable chunks will be deleted during runtime, and any other irrecoverable chunk located in the configured storage path directory will be deleted when Fluent-Bit starts.
	// +kubebuilder:validation:Enum:=on;off
	DeleteIrrecoverableChunks string `json:"deleteIrrecoverableChunks,omitempty"`
}

type Service struct {
	// If true go to background on start
	Daemon *bool `json:"daemon,omitempty"`
	// Interval to flush output
	FlushSeconds *int64 `json:"flushSeconds,omitempty"`
	// Wait time on exit
	GraceSeconds *int64 `json:"graceSeconds,omitempty"`
	// the error count to meet the unhealthy requirement, this is a sum for all output plugins in a defined HC_Period, example for output error: [2022/02/16 10:44:10] [ warn] [engine] failed to flush chunk '1-1645008245.491540684.flb', retry in 7 seconds: task_id=0, input=forward.1 > output=cloudwatch_logs.3 (out_id=3)
	// +kubebuilder:validation:Minimum:=1
	HcErrorsCount *int64 `json:"hcErrorsCount,omitempty"`
	// the retry failure count to meet the unhealthy requirement, this is a sum for all output plugins in a defined HC_Period, example for retry failure: [2022/02/16 20:11:36] [ warn] [engine] chunk '1-1645042288.260516436.flb' cannot be retried: task_id=0, input=tcp.3 > output=cloudwatch_logs.1
	// +kubebuilder:validation:Minimum:=1
	HcRetryFailureCount *int64 `json:"hcRetryFailureCount,omitempty"`
	// The time period by second to count the error and retry failure data point
	// +kubebuilder:validation:Minimum:=1
	HcPeriod *int64 `json:"hcPeriod,omitempty"`
	// enable Health check feature at http://127.0.0.1:2020/api/v1/health Note: Enabling this will not automatically configure kubernetes to use fluentbit's healthcheck endpoint
	HealthCheck *bool `json:"healthCheck,omitempty"`
	// Address to listen
	// +kubebuilder:validation:Pattern:="^\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}$"
	HttpListen string `json:"httpListen,omitempty"`
	// Port to listen
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	HttpPort *int32 `json:"httpPort,omitempty"`
	// If true enable statistics HTTP server
	HttpServer *bool `json:"httpServer,omitempty"`
	// File to log diagnostic output
	LogFile string `json:"logFile,omitempty"`
	// Diagnostic level (error/warning/info/debug/trace)
	// +kubebuilder:validation:Enum:=off;error;warning;info;debug;trace
	LogLevel string `json:"logLevel,omitempty"`
	// Optional 'parsers' config file (can be multiple)
	ParsersFile string `json:"parsersFile,omitempty"`
	// backward compatible
	ParsersFiles []string `json:"parsersFiles,omitempty"`
	// Configure a global environment for the storage layer in Service. It is recommended to configure the volume and volumeMount separately for this storage. The hostPath type should be used for that Volume in Fluentbit daemon set.
	Storage *Storage `json:"storage,omitempty"`
	// Per-namespace re-emitter configuration
	EmitterName        string `json:"emitterName,omitempty"`
	EmitterMemBufLimit string `json:"emitterMemBufLimit,omitempty"`
	EmitterStorageType string `json:"emitterStorageType,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=cfbc,scope=Cluster
// +genclient
// +genclient:nonNamespaced

// ClusterFluentBitConfig is the Schema for the cluster-level fluentbitconfigs API
type ClusterFluentBitConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FluentBitConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterFluentBitConfigList contains a list of ClusterFluentBitConfig
type ClusterFluentBitConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterFluentBitConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterFluentBitConfig{}, &ClusterFluentBitConfigList{})
}

func (s *Service) Params() *params.KVs {
	m := params.NewKVs()
	if s.Daemon != nil {
		m.Insert("Daemon", fmt.Sprint(*s.Daemon))
	}
	if s.FlushSeconds != nil {
		m.Insert("Flush", fmt.Sprint(*s.FlushSeconds))
	}
	if s.GraceSeconds != nil {
		m.Insert("Grace", fmt.Sprint(*s.GraceSeconds))
	}
	if s.HcErrorsCount != nil {
		m.Insert("HC_Errors_Count", fmt.Sprint(*s.HcErrorsCount))
	}
	if s.HcRetryFailureCount != nil {
		m.Insert("HC_Retry_Failure_Count", fmt.Sprint(*s.HcRetryFailureCount))
	}
	if s.HcPeriod != nil {
		m.Insert("HC_Period", fmt.Sprint(*s.HcPeriod))
	}
	if s.HealthCheck != nil {
		m.Insert("Health_Check", fmt.Sprint(*s.HealthCheck))
	}
	if s.HttpListen != "" {
		m.Insert("Http_Listen", s.HttpListen)
	}
	if s.HttpPort != nil {
		m.Insert("Http_Port", fmt.Sprint(*s.HttpPort))
	}
	if s.HttpServer != nil {
		m.Insert("Http_Server", fmt.Sprint(*s.HttpServer))
	}
	if s.LogFile != "" {
		m.Insert("Log_File", s.LogFile)
	}
	if s.LogLevel != "" {
		m.Insert("Log_Level", s.LogLevel)
	}
	if s.ParsersFile != "" {
		if s.ParsersFile == "parsers.conf" {
			// For backwards compatibility, if the "usual" parsers.conf is
			// configured, actually write the fully-qualified path  in order
			// to not break hot-reload.
			// See https://github.com/fluent/fluent-bit/issues/8275.
			m.Insert("Parsers_File", "/fluent-bit/etc/parsers.conf")
		} else {
			m.Insert("Parsers_File", s.ParsersFile)
		}
	}
	if len(s.ParsersFiles) != 0 {
		for _, parserFile := range s.ParsersFiles {
			m.Insert("Parsers_File", parserFile)
		}
	}
	if s.Storage != nil {
		if s.Storage.Path != "" {
			m.Insert("storage.path", s.Storage.Path)
		}
		if s.Storage.Sync != "" {
			m.Insert("storage.sync", s.Storage.Sync)
		}
		if s.Storage.Checksum != "" {
			m.Insert("storage.checksum", s.Storage.Checksum)
		}
		if s.Storage.BacklogMemLimit != "" {
			m.Insert("storage.backlog.mem_limit", s.Storage.BacklogMemLimit)
		}
		if s.Storage.Metrics != "" {
			m.Insert("storage.metrics", s.Storage.Metrics)
		}
		if s.Storage.MaxChunksUp != nil {
			m.Insert("storage.max_chunks_up", fmt.Sprint(*s.Storage.MaxChunksUp))
		}
		if s.Storage.DeleteIrrecoverableChunks != "" {
			m.Insert("storage.delete_irrecoverable_chunks", s.Storage.DeleteIrrecoverableChunks)
		}
	}
	return m
}

func (cfg ClusterFluentBitConfig) RenderMainConfig(
	sl plugins.SecretLoader, inputs ClusterInputList, filters ClusterFilterList,
	outputs ClusterOutputList, nsFilterLists []FilterList, nsOutputLists []OutputList, rewriteTagConfigs []string,
) (string, error) {
	var buf bytes.Buffer

	// The Service defines the global behaviour of the Fluent Bit engine.
	if cfg.Spec.Service != nil {
		buf.WriteString("[Service]\n")
		buf.WriteString(cfg.Spec.Service.Params().String())
	}

	inputSections, err := inputs.Load(sl)
	if err != nil {
		return "", err
	}

	filterSections, err := filters.Load(sl)
	if err != nil {
		return "", err
	}

	var nsFilterSections []string
	for _, nsFilterList := range nsFilterLists {
		if len(nsFilterList.Items) == 0 {
			continue
		}
		if nsFilterList.Items != nil {
			ns := nsFilterList.Items[0].Namespace
			namespacedSl := plugins.NewSecretLoader(sl.Client, ns)
			filters, err := nsFilterList.Load(namespacedSl)
			if err != nil {
				return "", err
			}
			nsFilterSections = append(nsFilterSections, filters)
		}
	}

	outputSections, err := outputs.Load(sl)
	if err != nil {
		return "", err
	}
	var nsOutputSections []string
	for _, nsOutputList := range nsOutputLists {
		if len(nsOutputList.Items) == 0 {
			continue
		}
		// The lists are per namespace, so get the namespace from the first item in a list
		if nsOutputList.Items != nil {
			ns := nsOutputList.Items[0].Namespace
			namespacedSl := plugins.NewSecretLoader(sl.Client, ns)
			outputs, err := nsOutputList.Load(namespacedSl)
			if err != nil {
				return "", err
			}
			nsOutputSections = append(nsOutputSections, outputs)
		}
	}

	if inputSections != "" && outputSections == "" && nsOutputSections == nil {
		outputSections = `[Output]
    Name    null
    Match   *`
	}

	buf.WriteString(inputSections)
	buf.WriteString(filterSections)
	for _, rtc := range rewriteTagConfigs {
		buf.WriteString(rtc)
	}
	for _, filters := range nsFilterSections {
		buf.WriteString(filters)
	}
	for _, outputs := range nsOutputSections {
		buf.WriteString(outputs)
	}
	buf.WriteString(outputSections)

	return buf.String(), nil
}

func (cfg ClusterFluentBitConfig) RenderParserConfig(
	sl plugins.SecretLoader, parsers ClusterParserList, nsParserLists []ParserList,
	nsClusterParserLists []ClusterParserList,
) (string, error) {
	var buf bytes.Buffer
	existingParsers := make(map[string]bool)
	parserSections, err := parsers.Load(sl, existingParsers)
	if err != nil {
		return "", err
	}

	buf.WriteString(parserSections)

	for _, parserListPerNS := range nsParserLists {
		if len(parserListPerNS.Items) == 0 {
			continue
		}
		if parserListPerNS.Items != nil {
			ns := parserListPerNS.Items[0].Namespace
			namespacedSl := plugins.NewSecretLoader(sl.Client, ns)
			nsParserSections, err := parserListPerNS.Load(namespacedSl)
			if err != nil {
				return "", err
			}
			buf.WriteString(nsParserSections)
		}
	}

	for _, item := range nsClusterParserLists {
		nsClusterParserSections, err := item.Load(sl, existingParsers)
		if err != nil {
			return "", err
		}
		buf.WriteString(nsClusterParserSections)
	}

	return buf.String(), nil
}

func (cfg ClusterFluentBitConfig) RenderMultilineParserConfig(
	sl plugins.SecretLoader, multilineParsers ClusterMultilineParserList, nsMultilineParserLists []MultilineParserList,
	nsClusterMultilineParserLists []ClusterMultilineParserList,
) (string, error) {
	var buf bytes.Buffer

	multilineParserSection, err := multilineParsers.Load(sl)
	if err != nil {
		return "", err
	}

	buf.WriteString(multilineParserSection)

	for _, nsmp := range nsMultilineParserLists {
		if len(nsmp.Items) == 0 {
			continue
		}
		if nsmp.Items != nil {
			ns := nsmp.Items[0].Namespace
			namespacedSl := plugins.NewSecretLoader(sl.Client, ns)
			nsmpSection, err := nsmp.Load(namespacedSl)
			if err != nil {
				return "", err
			}
			buf.WriteString(nsmpSection)
		}
	}

	for _, nscmp := range nsClusterMultilineParserLists {
		nscmpSection, err := nscmp.Load(sl)
		if err != nil {
			return "", err
		}
		buf.WriteString(nscmpSection)
	}

	return buf.String(), nil
}

// +kubebuilder:object:generate:=false

type Script struct {
	Name    string
	Content string
}

// +kubebuilder:object:generate:=false

// ByName implements sort.Interface for []Script based on the Name field.
type ByName []Script

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (cfg ClusterFluentBitConfig) RenderLuaScript(
	cl plugins.ConfigMapLoader, filters ClusterFilterList, namespace string,
) ([]Script, error) {

	scripts := make([]Script, 0)
	for _, f := range filters.Items {
		for _, p := range f.Spec.FilterItems {
			if p.Lua != nil {
				script, err := cl.LoadConfigMap(p.Lua.Script, namespace)
				if err != nil {
					return nil, err
				}
				scripts = append(scripts, Script{Name: p.Lua.Script.Key, Content: script})
			}
		}
	}

	sort.Sort(ByName(scripts))

	return scripts, nil
}
