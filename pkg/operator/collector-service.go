package operator

import (
	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	CollecotrMetricsPortName = "metrics"
	CollecotrMetricsPort     = 2020
	CollecotrTCPProtocolName = "TCP"
)

func MakeCollecotrService(co fluentbitv1alpha2.Collector) corev1.Service {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      co.Name,
			Namespace: co.Namespace,
			Labels:    co.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: co.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       CollecotrMetricsPortName,
					Port:       CollecotrMetricsPort,
					Protocol:   CollecotrTCPProtocolName,
					TargetPort: intstr.FromInt(CollecotrMetricsPort),
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

	return svc
}
