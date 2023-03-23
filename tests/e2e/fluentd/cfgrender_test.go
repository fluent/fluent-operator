package fluentd

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	fluentdtestcases "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/tests"
	"github.com/fluent/fluent-operator/v2/tests/utils"
)

const (
	mainAppKey = "app.conf"
)

var (
	ctx context.Context
)

// Run Test cases
var _ = Describe("Apply the fluentd forward CRs, comparing with the genrated configuraion is expected.", func() {

	ctx = context.TODO()
	BeforeEach(func() {
		time.Sleep(time.Second * 1)
	})

	AfterEach(func() {
		time.Sleep(time.Second * 1)
	})

	Describe("Test create and delete the fluentd CRs - 1", func() {
		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd clusterconfig and output with buffer example", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig1,
				&fluentdtestcases.FluentdClusterFilter1,
				&fluentdtestcases.FluentdClusterOutputBuffer,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdClusterCfgOutputWithBuffer)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 2", func() {
		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd clusterconfig and output to es", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig1,
				&fluentdtestcases.FluentdclusterOutput2ES,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdClusterCfgOutputES)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 3", func() {
		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd clusterconfig and output to kafka", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig1,
				&fluentdtestcases.FluentdClusterFilter1,
				&fluentdtestcases.FluentdClusterOutput2kafka,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdClusterCfgOutputKafka)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 4", func() {

		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd mixed configs and output to multi tenant", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig2,
				&fluentdtestcases.FluentdClusterOutputCluster,
				&fluentdtestcases.FluentdConfigUser1,
				&fluentdtestcases.FluentdClusterOutputLogOperator,
				&fluentdtestcases.FluentdOutputUser1,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdMixedCfgsMultiTenant)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 5", func() {

		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd mixed configs and output to es", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig1,
				&fluentdtestcases.FluentdConfig1,
				&fluentdtestcases.FluentdclusterOutput2ES,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdMixedCfgsOutputES)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 6", func() {

		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd namespaced config and output to es", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdConfig1,
				&fluentdtestcases.FluentdclusterOutput2ES,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdNamespacedCfgOutputES)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Test create and delete the fluentd CRs - 7", func() {
		It("E2E_FLUENTD_MAIN_APP_CONFIGURATION: fluentd clusterconfig and custom output to os", func() {

			objects := []client.Object{
				&fluentdtestcases.Fluentd,
				&fluentdtestcases.FluentdClusterFluentdConfig1,
				&fluentdtestcases.FluentdClusterOutputCustom,
			}

			err := CreateObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 2)

			seckey := types.NamespacedName{
				Namespace: fluentdtestcases.Fluentd.Namespace,
				Name:      fmt.Sprintf("%s-config", fluentdtestcases.Fluentd.Name),
			}
			config, err := GetCfgFromSecret(ctx, seckey)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(utils.ExpectedFluentdClusterCfgOutputCustom)).To(Equal(config))

			err = DeleteObjs(ctx, objects)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})

// CreateObjs create objs if not exists
func CreateObjs(ctx context.Context, objs []client.Object) error {
	for _, obj := range objs {
		err := k8sClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
		if err != nil {
			if errors.IsNotFound(err) {
				if obj.GetResourceVersion() != "" {
					obj.SetResourceVersion("")
				}
				if err := k8sClient.Create(ctx, obj); err != nil {
					return err
				}
				continue
			}

			return err
		}
	}

	return nil
}

// DeleteObjs delete objs with k8s client
func DeleteObjs(ctx context.Context, objs []client.Object) error {
	for _, obj := range objs {
		err := k8sClient.Delete(ctx, obj, client.GracePeriodSeconds(0))
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCfgFromSecret gets the configuration from the secret which mounted to the fluentd
func GetCfgFromSecret(ctx context.Context, key client.ObjectKey) (string, error) {
	var se corev1.Secret
	if err := k8sClient.Get(ctx, key, &se); err != nil {
		return "", err
	}

	data, ok := se.Data[mainAppKey]
	if !ok {
		return "", fmt.Errorf("Cannot get the main app conf")
	}

	return string(data), nil
}
