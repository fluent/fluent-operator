package operator

import (
	fluentbitv1alpha2 "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	FluentBitMetricsPortName = "metrics"
	FluentBitMetricsPort     = 2020
	FluentBitTCPProtocolName = "TCP"
)

func MakeFluentbitService(fb fluentbitv1alpha2.FluentBit) corev1.Service {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
			Labels:    fb.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: fb.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       FluentBitMetricsPortName,
					Port:       FluentBitMetricsPort,
					Protocol:   FluentBitTCPProtocolName,
					TargetPort: intstr.FromInt(FluentBitMetricsPort),
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

	return svc
}
