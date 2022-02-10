package operator

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "fluent.io/fluent-operator/apis/fluentd/v1alpha1"
)

const (
	ForwardPortName = "forward"
	HttpPortName    = "http"
)

func MakeFluentdService(fd fluentdv1alpha1.Fluentd) corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/name":      fd.Name,
		"app.kubernetes.io/instance":  "fluentd",
		"app.kubernetes.io/component": "fluentd",
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	// read inputs definition from globalInputs
	globalInputs := fd.Spec.GlobalInputs
	for _, input := range globalInputs {

		if input.Forward != nil {
			forwardPort := *input.Forward.Port
			if forwardPort == 0 {
				forwardPort = DefaultForwardPort
			}

			forwardContainerPort := corev1.ServicePort{
				Name:       DefaultForwardName,
				Port:       forwardPort,
				TargetPort: intstr.FromString(ForwardPortName),
				Protocol:   corev1.ProtocolTCP,
			}
			svc.Spec.Ports = append(svc.Spec.Ports, forwardContainerPort)
			continue
		}

		if input.Http != nil {
			httpPort := *input.Http.Port
			if httpPort == 0 {
				httpPort = DefaultHttpPort
			}
			httpContainerPort := corev1.ServicePort{
				Name:       DefaultHttpName,
				Port:       httpPort,
				TargetPort: intstr.FromString(HttpPortName),
				Protocol:   corev1.ProtocolTCP,
			}
			svc.Spec.Ports = append(svc.Spec.Ports, httpContainerPort)
		}
	}

	return svc
}
