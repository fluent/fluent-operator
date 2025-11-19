package fluentd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
)

var k8sClient client.Client

// Function to run the Ginkgo Test
func TestCompareFluentdMainAppConfig(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func(ctx context.Context) {
		path := os.Getenv("TESTCONFIG")
		if path == "" {
			path = fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		}

		cfg, err := clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			klog.Errorf("Failed to build config, err: %v", err)
			return
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())

		err = fluentdv1alpha1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		kc, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		k8sClient = kc
		Expect(k8sClient).NotTo(BeNil())

		_, _ = fmt.Fprintf(GinkgoWriter, "%s: Info: Setup Suite Execution\n", time.Now().Format(time.StampMilli))
	}, NodeTimeout(60*time.Second))

	AfterSuite(func() {
		By("After Suite Execution")
	})

	RunSpecs(t, "Test Fluentd Main Config Generated!")
}
