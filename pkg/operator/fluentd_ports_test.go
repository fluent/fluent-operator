package operator

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/input"
)

// spec.ports opens ports for inputs the operator can't derive a port from,
// such as a syslog source in a customPlugin, on the container and the Service.
func TestFluentdAdditionalPorts(t *testing.T) {
	syslog := corev1.ContainerPort{Name: "syslog", ContainerPort: 5140, Protocol: corev1.ProtocolUDP}
	syslogService := corev1.ServicePort{
		Name:       "syslog",
		Port:       5140,
		Protocol:   corev1.ProtocolUDP,
		TargetPort: intstr.FromInt32(5140),
	}

	for _, tc := range []struct {
		name               string
		inputs             []input.Input
		ports              []corev1.ContainerPort
		wantContainerPorts []string
		wantServicePorts   []corev1.ServicePort
	}{
		{
			name:               "only additional ports",
			ports:              []corev1.ContainerPort{syslog},
			wantContainerPorts: []string{MetricsName, "syslog"},
			wantServicePorts:   []corev1.ServicePort{syslogService},
		},
		{
			name:               "forward input and additional ports",
			inputs:             []input.Input{{Forward: &input.Forward{Port: new(int32(24224))}}},
			ports:              []corev1.ContainerPort{syslog},
			wantContainerPorts: []string{MetricsName, DefaultForwardName, "syslog"},
			wantServicePorts: []corev1.ServicePort{
				{
					Name:       DefaultForwardName,
					Port:       24224,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString(FluentdForwardPortName),
				},
				syslogService,
			},
		},
		{
			name:               "no additional ports",
			inputs:             []input.Input{{Forward: &input.Forward{Port: new(int32(24224))}}},
			wantContainerPorts: []string{MetricsName, DefaultForwardName},
			wantServicePorts: []corev1.ServicePort{
				{
					Name:       DefaultForwardName,
					Port:       24224,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString(FluentdForwardPortName),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fd := fluentdv1alpha1.Fluentd{
				ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging"},
				Spec:       fluentdv1alpha1.FluentdSpec{GlobalInputs: tc.inputs, Ports: tc.ports},
			}

			containerPorts := makeFluentdPorts(fd)
			if len(containerPorts) != len(tc.wantContainerPorts) {
				t.Fatalf("container ports = %+v, want names %v", containerPorts, tc.wantContainerPorts)
			}
			for i, name := range tc.wantContainerPorts {
				if containerPorts[i].Name != name {
					t.Fatalf("container port %d = %+v, want name %q", i, containerPorts[i], name)
				}
			}
			if len(tc.ports) > 0 && containerPorts[len(containerPorts)-1] != tc.ports[len(tc.ports)-1] {
				t.Fatalf("last container port = %+v, want %+v", containerPorts[len(containerPorts)-1], tc.ports[len(tc.ports)-1])
			}

			svc := MakeFluentdService(fd)
			if len(svc.Spec.Ports) != len(tc.wantServicePorts) {
				t.Fatalf("service ports = %+v, want %+v", svc.Spec.Ports, tc.wantServicePorts)
			}
			for i := range tc.wantServicePorts {
				if svc.Spec.Ports[i] != tc.wantServicePorts[i] {
					t.Fatalf("service port %d = %+v, want %+v", i, svc.Spec.Ports[i], tc.wantServicePorts[i])
				}
			}
		})
	}
}

// Several inputs of one kind each get their own port; the names have to be
// unique or the StatefulSet and the Service are rejected.
func TestFluentdMultipleInputsOfOneKind(t *testing.T) {
	fd := fluentdv1alpha1.Fluentd{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging"},
		Spec: fluentdv1alpha1.FluentdSpec{GlobalInputs: []input.Input{
			{Forward: &input.Forward{Port: new(int32(24224))}},
			{Forward: &input.Forward{Port: new(int32(25224))}},
			{Http: &input.Http{Port: new(int32(9880))}},
			{Http: &input.Http{Port: new(int32(9881))}},
			{Forward: &input.Forward{Port: new(int32(26224))}},
		}},
	}

	wantNames := []string{MetricsName, "forward", "forward-2", "http", "http-2", "forward-3"}
	containerPorts := makeFluentdPorts(fd)
	if len(containerPorts) != len(wantNames) {
		t.Fatalf("container ports = %+v, want names %v", containerPorts, wantNames)
	}
	for i, name := range wantNames {
		if containerPorts[i].Name != name {
			t.Fatalf("container port %d = %+v, want name %q", i, containerPorts[i], name)
		}
	}

	servicePorts := MakeFluentdService(fd).Spec.Ports
	if len(servicePorts) != len(wantNames)-1 {
		t.Fatalf("service ports = %+v, want %d", servicePorts, len(wantNames)-1)
	}
	for i, svcPort := range servicePorts {
		containerPort := containerPorts[i+1]
		if svcPort.Name != containerPort.Name || svcPort.TargetPort != intstr.FromString(containerPort.Name) || svcPort.Port != containerPort.ContainerPort {
			t.Fatalf("service port %d = %+v does not match container port %+v", i, svcPort, containerPort)
		}
	}
}
