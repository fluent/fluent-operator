package operator

import (
	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	CollectorMetricsPortName = "metrics"
	CollectorMetricsPort     = 2020
	CollectorTCPProtocolName = "TCP"
)

func MakeCollectorService(co fluentbitv1alpha2.Collector) *corev1.Service {
	var name string
	var labels map[string]string

	if co.Spec.Service.Name != "" {
		name = co.Spec.Service.Name
	} else {
		name = co.Name
	}

	if len(co.Spec.Service.Labels) != 0 {
		labels = co.Spec.Service.Labels
	} else {
		labels = co.Labels
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: co.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: co.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       CollectorMetricsPortName,
					Port:       CollectorMetricsPort,
					Protocol:   CollectorTCPProtocolName,
					TargetPort: intstr.FromInt(CollectorMetricsPort),
				},
			},
		},
	}

	if co.Spec.Ports != nil {
		for _, port := range co.Spec.Ports {
			svc.Spec.Ports = append(svc.Spec.Ports, corev1.ServicePort{
				Name:       port.Name,
				Port:       port.ContainerPort,
				Protocol:   port.Protocol,
				TargetPort: intstr.FromInt(int(port.ContainerPort)),
			})
		}
	}

	if len(co.Spec.Service.Annotations) != 0 {
		svc.ObjectMeta.Annotations = co.Spec.Service.Annotations
	}

	return &svc
}
