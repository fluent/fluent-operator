package fluentd

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v3/tests/utils"
)

var (
	// Fluentd instance for label selector tests
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

	// Grep Filter
	FilterGrepRaw = `
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

			// Parse YAML into objects
			var fluentd fluentdv1alpha1.Fluentd
			err := yaml.Unmarshal([]byte(FluentdLabelSelectorRaw), &fluentd)
			Expect(err).NotTo(HaveOccurred())

			var fluentdConfig fluentdv1alpha1.FluentdConfig
			err = yaml.Unmarshal([]byte(FluentdConfigLabelSelectorRaw), &fluentdConfig)
			Expect(err).NotTo(HaveOccurred())

			var testFilter fluentdv1alpha1.Filter
			err = yaml.Unmarshal([]byte(FilterLabelSelectorRaw), &testFilter)
			Expect(err).NotTo(HaveOccurred())

			var testOutput fluentdv1alpha1.Output
			err = yaml.Unmarshal([]byte(OutputLabelSelectorRaw), &testOutput)
			Expect(err).NotTo(HaveOccurred())

			// Create all objects
			objects := []client.Object{
				&fluentd,
				&fluentdConfig,
				&testFilter,
				&testOutput,
			}

			err = CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			// Wait for reconciliation
			time.Sleep(time.Second * 3)

			// Get the generated configuration
			seckey := types.NamespacedName{
				Namespace: fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			// Verify that the filter configuration is present
			// Before the fix, the filter would not be loaded because it was written to the inputs list
			Expect(string(utils.ExpectedFluentdNamespacedCfgFilterOutputSelector)).To(Equal(config))

			// Clean up
			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})

		It("E2E_FLUENTD_NAMESPACE_MIXED_SELECTORS: FluentdConfig with only filterSelector", func() {

			// Parse YAML into objects
			var fluentd fluentdv1alpha1.Fluentd
			err := yaml.Unmarshal([]byte(FluentdLabelSelectorRaw), &fluentd)
			Expect(err).NotTo(HaveOccurred())

			var fluentdConfig fluentdv1alpha1.FluentdConfig
			err = yaml.Unmarshal([]byte(FluentdConfigFilterOnlyRaw), &fluentdConfig)
			Expect(err).NotTo(HaveOccurred())

			var grepFilter fluentdv1alpha1.Filter
			err = yaml.Unmarshal([]byte(FilterGrepRaw), &grepFilter)
			Expect(err).NotTo(HaveOccurred())

			var clusterOutput fluentdv1alpha1.ClusterOutput
			err = yaml.Unmarshal([]byte(ClusterOutputStdoutRaw), &clusterOutput)
			Expect(err).NotTo(HaveOccurred())

			objects := []client.Object{
				&fluentd,
				&fluentdConfig,
				&grepFilter,
				&clusterOutput,
			}

			err = CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 3)

			seckey := types.NamespacedName{
				Namespace: fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(string(utils.ExpectedFluentdNamespacedCfgFilterSelector)).To(Equal(config))

			// Clean up
			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
