package operator

import (
	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	FluentBitMetricsPortName = "metrics"
	FluentBitTCPProtocolName = "TCP"
)

func MakeFluentbitService(fb fluentbitv1alpha2.FluentBit) *corev1.Service {
	var name string
	var labels map[string]string

	if fb.Spec.Service.Name != "" {
		name = fb.Spec.Service.Name
	} else {
		name = fb.Name
	}

	if len(fb.Spec.Service.Labels) != 0 {
		labels = fb.Spec.Service.Labels
	} else {
		labels = fb.Labels
	}

	var FluentBitMetricsPort int32
	if fb.Spec.MetricsPort != 0 {
		FluentBitMetricsPort = fb.Spec.MetricsPort
	} else {
		FluentBitMetricsPort = 2020
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: fb.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: fb.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       FluentBitMetricsPortName,
					Port:       FluentBitMetricsPort,
					Protocol:   FluentBitTCPProtocolName,
					TargetPort: intstr.FromInt(int(FluentBitMetricsPort)),
				},
			},
		},
	}

	if fb.Spec.Ports != nil {
		for _, port := range fb.Spec.Ports {
			svc.Spec.Ports = append(svc.Spec.Ports, corev1.ServicePort{
				Name:       port.Name,
				Port:       port.ContainerPort,
				Protocol:   port.Protocol,
				TargetPort: intstr.FromInt(int(port.ContainerPort)),
			})
		}
	}

	if len(fb.Spec.Service.Annotations) != 0 {
		svc.ObjectMeta.Annotations = fb.Spec.Service.Annotations
	}

	return &svc
}
