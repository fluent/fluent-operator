package utils

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once

	ExpectedFluentdClusterCfgOutputES                []byte
	ExpectedFluentdClusterCfgOutputKafka             []byte
	ExpectedFluentdClusterCfgOutputWithBuffer        []byte
	ExpectedFluentdMixedCfgsMultiTenant              []byte
	ExpectedFluentdMixedCfgsOutputES                 []byte
	ExpectedFluentdNamespacedCfgOutputES             []byte
	ExpectedDuplicateRemovalCRSPECS                  []byte
	ExpectedFluentdClusterCfgOutputCustom            []byte
	ExpectedFluentdNamespacedCfgFilterOutputSelector []byte
	ExpectedFluentdNamespacedCfgFilterSelector       []byte
)

func init() {
	once.Do(func() {
		ExpectedFluentdClusterCfgOutputES =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-cluster-cfg-output-es.cfg")
		ExpectedFluentdClusterCfgOutputKafka =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-cluster-cfg-output-kafka.cfg")
		ExpectedFluentdClusterCfgOutputWithBuffer =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-cluster-cfg-output-buffer-example.cfg")
		ExpectedFluentdMixedCfgsMultiTenant =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-mixed-cfgs-multi-tenant-output.cfg")
		ExpectedFluentdMixedCfgsOutputES =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-mixed-cfgs-output-es.cfg")
		ExpectedFluentdNamespacedCfgOutputES =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-namespaced-cfg-output-es.cfg")
		ExpectedDuplicateRemovalCRSPECS =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/duplicate-removal-cr-specs.cfg")
		ExpectedFluentdClusterCfgOutputCustom =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-cluster-cfg-output-custom.cfg")
		ExpectedFluentdNamespacedCfgFilterOutputSelector =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-namespaced-cfg-filter-output-selector.cfg")
		ExpectedFluentdNamespacedCfgFilterSelector =
			getExpectedCfg("./apis/fluentd/v1alpha1/tests/expected/fluentd-namespaced-cfg-filter-selector.cfg")
	})
}

func getExpectedCfg(path string) []byte {
	if !filepath.IsAbs(path) {
		var projPath string
		pwd, _ := os.Getwd()
		if _, file := filepath.Split(pwd); file == "fluentd" {
			// For debug mode
			projPath = filepath.Dir(filepath.Dir(filepath.Dir(pwd)))
		} else {
			projPath = pwd
		}
		path = filepath.Join(projPath, path)
	}
	data, _ := os.ReadFile(path)
	return data
}
