package operator

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
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

	serviceType := corev1.ServiceTypeClusterIP
	if fd.Spec.Service.Type != nil {
		serviceType = *fd.Spec.Service.Type
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: fd.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     serviceType,
		},
	}

	// Read inputs definition from globalInputs
	globalInputs := fd.Spec.GlobalInputs
	firstForwardPort := true
	firstHttpPort := true
	for _, input := range globalInputs {

		if input.Forward != nil {
			forwardPort := *input.Forward.Port
			if forwardPort == 0 {
				forwardPort = DefaultForwardPort
			}

			forwardName := DefaultForwardName
			if !firstForwardPort {
				forwardName = fmt.Sprintf("%s-%d", DefaultForwardName, forwardPort)
			}
			forwardContainerPort := corev1.ServicePort{
				Name:       forwardName,
				Port:       forwardPort,
				TargetPort: intstr.FromString(forwardName),
				Protocol:   corev1.ProtocolTCP,
			}
			svc.Spec.Ports = append(svc.Spec.Ports, forwardContainerPort)
			firstForwardPort = false
			continue
		}

		if input.Http != nil {
			httpPort := *input.Http.Port
			if httpPort == 0 {
				httpPort = DefaultHttpPort
			}
			httpName := DefaultHttpName
			if !firstHttpPort {
				httpName = fmt.Sprintf("%s-%d", DefaultHttpName, httpPort)
			}
			httpContainerPort := corev1.ServicePort{
				Name:       httpName,
				Port:       httpPort,
				TargetPort: intstr.FromString(httpName),
				Protocol:   corev1.ProtocolTCP,
			}
			svc.Spec.Ports = append(svc.Spec.Ports, httpContainerPort)
			firstHttpPort = false
		}
	}

	if len(fd.Spec.Service.Annotations) != 0 {
		svc.Annotations = fd.Spec.Service.Annotations
	}

	return &svc
}
