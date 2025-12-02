package fluentd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	cfgrender "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/tests"
	"github.com/fluent/fluent-operator/v3/tests/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	once sync.Once

	// Fluentd instance for label selector tests
	Fluentd                 fluentdv1alpha1.Fluentd
	FluentdLabelSelectorRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Fluentd
metadata:
  name: fluentd-label-selector-test
  namespace: fluent
  labels:
    app.kubernetes.io/name: fluentd
spec:
  globalInputs:
  - forward:
      bind: 0.0.0.0
      port: 24224
  replicas: 1
  image: ghcr.io/fluent/fluent-operator/fluentd:v1.19.1
  fluentdCfgSelector:
    matchLabels:
      label.config.fluentd.fluent.io/enabled: "true"
`

	// FluentdConfig with filterSelector and outputSelector
	FluentdConfig                 fluentdv1alpha1.FluentdConfig
	FluentdConfigLabelSelectorRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: FluentdConfig
metadata:
  name: fluentd-config-label-selector-test
  namespace: fluent
  labels:
    label.config.fluentd.fluent.io/enabled: "true"
spec:
  filterSelector:
    matchLabels:
      filter.fluentd.fluent.io/enabled: "true"
      filter.fluentd.fluent.io/mode: "namespace"
  outputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
      output.fluentd.fluent.io/mode: "namespace"
`

	// Filter with matching labels
	FluentdFilter          fluentdv1alpha1.Filter
	FilterLabelSelectorRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Filter
metadata:
  name: test-filter-recordtransformer
  namespace: fluent
  labels:
    filter.fluentd.fluent.io/enabled: "true"
    filter.fluentd.fluent.io/mode: "namespace"
spec:
  filters:
  - recordTransformer:
      records:
      - key: hostname
        value: test-host
      - key: environment
        value: testing
`

	// Output with matching labels
	FluentdOutput          fluentdv1alpha1.Output
	OutputLabelSelectorRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: test-output-stdout
  namespace: fluent
  labels:
    output.fluentd.fluent.io/enabled: "true"
    output.fluentd.fluent.io/mode: "namespace"
spec:
  outputs:
  - stdout: {}
`

	// FluentdConfig with only filterSelector
	FluentdConfigFilterOnly    fluentdv1alpha1.FluentdConfig
	FluentdConfigFilterOnlyRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: FluentdConfig
metadata:
  name: fluentd-config-filter-only
  namespace: fluent
  labels:
    label.config.fluentd.fluent.io/enabled: "true"
spec:
  filterSelector:
    matchLabels:
      filter.fluentd.fluent.io/type: "grep"
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
`

	// Grep Filter with matching labels
	FluentdFilterGrep fluentdv1alpha1.Filter
	FilterGrepRaw     = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Filter
metadata:
  name: test-filter-grep
  namespace: fluent
  labels:
    filter.fluentd.fluent.io/type: "grep"
spec:
  filters:
  - grep:
      regexp:
      - key: level
        pattern: /error/
`

	// ClusterOutput for filter-only test
	FluentdClusterOutput   fluentdv1alpha1.ClusterOutput
	ClusterOutputStdoutRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: cluster-output-stdout-filter-test
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec:
  outputs:
  - stdout: {}
`
)

func init() {
	once.Do(setupFluentdObjects)
}

func setupFluentdObjects() {
	cfgrender.MustParseIntoObject(FluentdLabelSelectorRaw, &Fluentd)
	cfgrender.MustParseIntoObject(FluentdConfigLabelSelectorRaw, &FluentdConfig)
	cfgrender.MustParseIntoObject(FluentdConfigFilterOnlyRaw, &FluentdConfigFilterOnly)
	cfgrender.MustParseIntoObject(FilterLabelSelectorRaw, &FluentdFilter)
	cfgrender.MustParseIntoObject(OutputLabelSelectorRaw, &FluentdOutput)
	cfgrender.MustParseIntoObject(FilterGrepRaw, &FluentdFilterGrep)
	cfgrender.MustParseIntoObject(ClusterOutputStdoutRaw, &FluentdClusterOutput)
}

// Helper function to run a fluentd label selector test
func testFluentdLabelSelector(
	ctx context.Context,
	expectedConfig []byte,
	fluentd fluentdv1alpha1.Fluentd,
	objects []client.Object,
) {
	// Create all objects
	err := CreateObjs(ctx, objects)
	Expect(err).NotTo(HaveOccurred())

	// Ensure cleanup runs even if the test fails
	DeferCleanup(func() {
		// Clean up all objects
		err := DeleteObjs(ctx, objects)
		if err != nil {
			// Log the error but don't fail the cleanup
			fmt.Printf("Warning: failed to cleanup objects: %v\n", err)
		}
	})

	// Wait for reconciliation
	time.Sleep(time.Second * 3)

	// Get the generated configuration
	seckey := types.NamespacedName{
		Namespace: fluentd.Namespace,
		Name:      fmt.Sprintf("%s-config", fluentd.Name),
	}
	config, err := GetCfgFromSecret(ctx, seckey)
	Expect(err).NotTo(HaveOccurred())

	// Verify that the configuration matches expected
	Expect(strings.TrimRight(string(expectedConfig), "\r\n")).To(Equal(config))
}

// This test verifies the fix for the bug where filterSelector and outputSelector
// were incorrectly writing to the inputs list instead of their respective lists.
var _ = Describe("Test FluentdConfig with namespace-level filter and output selectors", func() {

	ctx := context.TODO()

	BeforeEach(func() {
		time.Sleep(time.Second * 1)
	})

	AfterEach(func() {
		time.Sleep(time.Second * 1)
	})

	Describe("Test namespace-level resources with label selectors", func() {
		It("E2E_FLUENTD_NAMESPACE_FILTER_OUTPUT_SELECTORS: FluentdConfig with filterSelector and outputSelector", func() {
			testFluentdLabelSelector(
				ctx,
				utils.ExpectedFluentdNamespacedCfgFilterOutputSelector,
				Fluentd,
				[]client.Object{
					&Fluentd,
					&FluentdConfig,
					&FluentdFilter,
					&FluentdOutput,
				},
			)
		})

		It("E2E_FLUENTD_NAMESPACE_MIXED_SELECTORS: FluentdConfig with only filterSelector", func() {
			testFluentdLabelSelector(
				ctx,
				utils.ExpectedFluentdNamespacedCfgFilterSelector,
				Fluentd,
				[]client.Object{
					&Fluentd,
					&FluentdConfigFilterOnly,
					&FluentdFilterGrep,
					&FluentdClusterOutput,
				},
			)
		})
	})
})
