package operator

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
)

const (
	FluentdForwardPortName = "forward"
	FluentdHttpPortName    = "http"
)

func MakeFluentdService(fd fluentdv1alpha1.Fluentd) *corev1.Service {
	var name string
	var labels map[string]string

	if fd.Spec.Service.Name != "" {
		name = fd.Spec.Service.Name
	} else {
		name = fd.Name
	}

	if len(fd.Spec.Service.Labels) != 0 {
		labels = fd.Spec.Service.Labels
	} else {
		labels = map[string]string{
			"app.kubernetes.io/name":      name,
			"app.kubernetes.io/instance":  "fluentd",
			"app.kubernetes.io/component": "fluentd",
		}
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
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
				TargetPort: intstr.FromString(FluentdForwardPortName),
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
				TargetPort: intstr.FromString(FluentdHttpPortName),
				Protocol:   corev1.ProtocolTCP,
			}
			svc.Spec.Ports = append(svc.Spec.Ports, httpContainerPort)
		}
	}

	if len(fd.Spec.Service.Annotations) != 0 {
		svc.ObjectMeta.Annotations = fd.Spec.Service.Annotations
	}

	return &svc
}
