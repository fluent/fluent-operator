package fluentd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/input"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/output"
)

// generateRandomID creates a cryptographically random hex string
func generateRandomID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

var _ = Describe("Fluentd E2E Deployment Test", func() {
	var cancel context.CancelFunc
	var ctx context.Context

	BeforeEach(func() {
		// Create context with timeout to prevent hung tests
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
	})

	Describe("Deploying Fluentd CR", func() {
		var (
			fluentdCR     *fluentdv1alpha1.Fluentd
			fluentdConfig *fluentdv1alpha1.FluentdConfig
			clusterOutput *fluentdv1alpha1.ClusterOutput
			namespace     string
		)

		BeforeEach(func() {
			// Generate a unique namespace using crypto/rand for true isolation
			namespace = fmt.Sprintf("fluentd-e2e-%s", generateRandomID())
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			// Handle case where namespace might already exist from crashed previous run
			err := k8sClient.Create(ctx, ns)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				Fail(fmt.Sprintf("Failed to create namespace: %v", err))
			}

			// Create Fluentd CR with proper GlobalInputs type
			fluentdCR = &fluentdv1alpha1.Fluentd{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fluentd-instance",
					Namespace: namespace,
					Labels: map[string]string{
						"app.kubernetes.io/name":     "fluentd",
						"app.kubernetes.io/instance": "fluentd-instance",
					},
				},
				Spec: fluentdv1alpha1.FluentdSpec{
					Replicas: ptr.To(int32(1)),
					GlobalInputs: []input.Input{
						{
							Forward: &input.Forward{
								Bind: ptr.To("0.0.0.0"),
								Port: ptr.To(int32(24224)),
							},
						},
					},
					// Explicitly set image as operator doesn't provide a default yet
					Image: "ghcr.io/fluent/fluent-operator/fluentd:v1.19.1",
					// Use EmptyDir for buffers to avoid PVC provisioning issues in CI
					BufferVolume: &fluentdv1alpha1.BufferVolume{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
					FluentdCfgSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"config.fluentd.fluent.io/enabled": "true",
						},
					},
				},
			}

			// Create a ClusterOutput for stdout (minimal working config)
			clusterOutput = &fluentdv1alpha1.ClusterOutput{
				ObjectMeta: metav1.ObjectMeta{
					Name: "fluentd-output-stdout",
					Labels: map[string]string{
						"output.fluentd.fluent.io/enabled": "true",
					},
				},
				Spec: fluentdv1alpha1.ClusterOutputSpec{
					Outputs: []output.Output{
						{
							Stdout: &output.Stdout{},
						},
					},
				},
			}

			// Create FluentdConfig to wire everything together
			fluentdConfig = &fluentdv1alpha1.FluentdConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fluentd-config",
					Namespace: namespace,
					Labels: map[string]string{
						"config.fluentd.fluent.io/enabled": "true",
					},
				},
				Spec: fluentdv1alpha1.FluentdConfigSpec{
					ClusterOutputSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"output.fluentd.fluent.io/enabled": "true",
						},
					},
				},
			}

			DeferCleanup(func() {
				// Use a fresh context for cleanup to avoid timeout issues
				cleanupCtx := context.Background()

				// Delete all CRs (ignore NotFound errors for idempotency)
				if clusterOutput != nil {
					_ = client.IgnoreNotFound(k8sClient.Delete(cleanupCtx, clusterOutput))
				}
				if fluentdConfig != nil {
					_ = client.IgnoreNotFound(k8sClient.Delete(cleanupCtx, fluentdConfig))
				}
				if fluentdCR != nil {
					_ = client.IgnoreNotFound(k8sClient.Delete(cleanupCtx, fluentdCR))
				}

				// Wait for StatefulSet to be deleted (find by label, not name)
				if fluentdCR != nil {
					Eventually(func() bool {
						stsList := &appsv1.StatefulSetList{}
						err := k8sClient.List(cleanupCtx, stsList, client.InNamespace(namespace))
						if err != nil {
							return false
						}
						// StatefulSet should be gone
						return len(stsList.Items) == 0
					}, time.Minute, time.Second).Should(BeTrue(), "StatefulSet should be deleted")
				}

				// Delete namespace and wait for it to be gone
				ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
				_ = client.IgnoreNotFound(k8sClient.Delete(cleanupCtx, ns))
				Eventually(func() bool {
					err := k8sClient.Get(cleanupCtx, types.NamespacedName{Name: namespace}, &corev1.Namespace{})
					return apierrors.IsNotFound(err)
				}, 2*time.Minute, time.Second).Should(BeTrue(), "Namespace should be deleted")
			})
		})

		It("Should create a healthy Fluentd StatefulSet", func() {
			By("Creating the ClusterOutput")
			Expect(k8sClient.Create(ctx, clusterOutput)).To(Succeed())

			By("Creating the FluentdConfig")
			Expect(k8sClient.Create(ctx, fluentdConfig)).To(Succeed())

			By("Creating the Fluentd Custom Resource")
			Expect(k8sClient.Create(ctx, fluentdCR)).To(Succeed())

			By("Verifying StatefulSet creation and readiness")

			// Find StatefulSet by label instead of assuming name
			Eventually(func() bool {
				stsList := &appsv1.StatefulSetList{}
				err := k8sClient.List(ctx, stsList,
					client.InNamespace(namespace),
					client.MatchingLabels{"app.kubernetes.io/name": fluentdCR.Name})
				if err != nil || len(stsList.Items) == 0 {
					return false
				}
				return true
			}, time.Minute, time.Second).Should(BeTrue(), "StatefulSet should be created by the Operator")

			// Check for Ready Replicas (Real Workload Health)
			Eventually(func() int32 {
				// Refresh StatefulSet status
				stsList := &appsv1.StatefulSetList{}
				_ = k8sClient.List(ctx, stsList,
					client.InNamespace(namespace),
					client.MatchingLabels{"app.kubernetes.io/name": fluentdCR.Name})
				if len(stsList.Items) == 0 {
					return 0
				}
				return stsList.Items[0].Status.ReadyReplicas
			}, 5*time.Minute, 2*time.Second).Should(Equal(*fluentdCR.Spec.Replicas),
				"StatefulSet should have expected number of ready replicas")

			By("Verifying Pod Status and Container Readiness")
			Eventually(func() bool {
				podList := &corev1.PodList{}
				// List pods owned by the StatefulSet
				err := k8sClient.List(ctx, podList, client.InNamespace(namespace))
				if err != nil {
					return false
				}
				for _, pod := range podList.Items {
					// Check if pod is owned by a StatefulSet
					for _, owner := range pod.OwnerReferences {
						if owner.Kind == "StatefulSet" {
							// Verify pod is running
							if pod.Status.Phase != corev1.PodRunning {
								continue
							}

							// Check ALL containers are ready (not just pod condition)
							allContainersReady := true
							for _, cs := range pod.Status.ContainerStatuses {
								if !cs.Ready {
									allContainersReady = false
									break
								}
							}

							// Also check pod ready condition
							podReady := false
							for _, condition := range pod.Status.Conditions {
								if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
									podReady = true
									break
								}
							}

							if allContainersReady && podReady {
								return true
							}
						}
					}
				}
				return false
			}, 5*time.Minute, 2*time.Second).Should(BeTrue(),
				"At least one Fluentd Pod should be Running with all containers Ready")
		})
	})
})
