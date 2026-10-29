package operator

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/input"
)

// Inputs that leave their port unset listen on Fluentd's defaults, so the
// container ports and the Service must point there too.
func TestFluentdInputPorts(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    input.Input
		portName string
		want     int32
	}{
		{name: "forward without a port", input: input.Input{Forward: &input.Forward{}}, portName: DefaultForwardName, want: 24224},
		{
			name:     "forward with a port",
			input:    input.Input{Forward: &input.Forward{Port: int32Ptr(24225)}},
			portName: DefaultForwardName,
			want:     24225,
		},
		{name: "http without a port", input: input.Input{Http: &input.Http{}}, portName: DefaultHttpName, want: 9880},
		{
			name:     "http with a port",
			input:    input.Input{Http: &input.Http{Port: int32Ptr(9881)}},
			portName: DefaultHttpName,
			want:     9881,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fd := fluentdv1alpha1.Fluentd{
				ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging"},
				Spec:       fluentdv1alpha1.FluentdSpec{GlobalInputs: []input.Input{tc.input}},
			}

			var containerPort *corev1.ContainerPort
			for _, p := range makeFluentdPorts(fd) {
				if p.Name == tc.portName {
					containerPort = &p
				}
			}
			if containerPort == nil || containerPort.ContainerPort != tc.want {
				t.Fatalf("container port %q = %v, want %d", tc.portName, containerPort, tc.want)
			}

			svc := MakeFluentdService(fd)
			if len(svc.Spec.Ports) != 1 {
				t.Fatalf("service ports = %v, want one", svc.Spec.Ports)
			}
			if got := svc.Spec.Ports[0]; got.Port != tc.want || got.TargetPort != intstr.FromString(tc.portName) {
				t.Fatalf("service port = %+v, want port %d targeting %q", got, tc.want, tc.portName)
			}
		})
	}
}

func int32Ptr(v int32) *int32 { return &v }
