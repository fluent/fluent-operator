package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

const testLogPath = "/var/lib/docker/containers"

func secretVolume(name, secretName string) corev1.Volume {
	return corev1.Volume{Name: name, VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{SecretName: secretName},
	}}
}

func hostPathVolume(name, path string) corev1.Volume {
	return corev1.Volume{Name: name, VolumeSource: corev1.VolumeSource{
		HostPath: &corev1.HostPathVolumeSource{Path: path},
	}}
}

func TestMakeDaemonSetDefaultVolumes(t *testing.T) {
	g := NewWithT(t)
	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{FluentBitConfigName: "fb-config"},
	}

	ds := MakeDaemonSet(fb, testLogPath)
	podSpec := ds.Spec.Template.Spec
	none := corev1.MountPropagationNone

	g.Expect(podSpec.Volumes).To(Equal([]corev1.Volume{
		secretVolume("config", "fb-config"),
		hostPathVolume("varlibcontainers", testLogPath),
		hostPathVolume("varlogs", "/var/log"),
		hostPathVolume("systemd", "/var/log/journal"),
	}))
	g.Expect(podSpec.Containers[0].VolumeMounts).To(Equal([]corev1.VolumeMount{
		{Name: "config", ReadOnly: true, MountPath: "/fluent-bit/config"},
		{Name: "varlibcontainers", ReadOnly: true, MountPath: testLogPath, MountPropagation: &none},
		{Name: "varlogs", ReadOnly: true, MountPath: "/var/log/", MountPropagation: &none},
		{Name: "systemd", ReadOnly: true, MountPath: "/var/log/journal", MountPropagation: &none},
	}))
}

func TestMakeDaemonSetInternalMountPropagation(t *testing.T) {
	g := NewWithT(t)
	hostToContainer := corev1.MountPropagationHostToContainer
	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{InternalMountPropagation: &hostToContainer},
	}

	for _, m := range MakeDaemonSet(fb, testLogPath).Spec.Template.Spec.Containers[0].VolumeMounts {
		if m.Name == "config" {
			g.Expect(m.MountPropagation).To(BeNil())
			continue
		}
		g.Expect(m.MountPropagation).To(HaveValue(Equal(hostToContainer)), m.Name)
	}
}

func TestMakeDaemonSetCustomVolumes(t *testing.T) {
	customVolumes := []corev1.Volume{
		{Name: "pvc-volume", VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "pvc"},
		}},
		{Name: "emptydir-volume", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	customVolumeMounts := []corev1.VolumeMount{
		{Name: "pvc-volume", MountPath: "/mnt/pvc"},
		{Name: "emptydir-volume", MountPath: "/mnt/emptydir"},
	}

	for _, disableLogVolumes := range []bool{false, true} {
		g := NewWithT(t)
		fb := fluentbitv1alpha2.FluentBit{
			Spec: fluentbitv1alpha2.FluentBitSpec{
				DisableLogVolumes: disableLogVolumes,
				Volumes:           customVolumes,
				VolumesMounts:     customVolumeMounts,
			},
		}

		podSpec := MakeDaemonSet(fb, testLogPath).Spec.Template.Spec
		volumeNames := make([]string, 0, len(podSpec.Volumes))
		for _, v := range podSpec.Volumes {
			volumeNames = append(volumeNames, v.Name)
		}
		mountNames := make([]string, 0, len(podSpec.Containers[0].VolumeMounts))
		for _, m := range podSpec.Containers[0].VolumeMounts {
			mountNames = append(mountNames, m.Name)
		}

		want := []string{"config", "varlibcontainers", "varlogs", "systemd", "pvc-volume", "emptydir-volume"}
		if disableLogVolumes {
			want = []string{"config", "pvc-volume", "emptydir-volume"}
		}
		g.Expect(volumeNames).To(Equal(want))
		g.Expect(mountNames).To(Equal(want))
	}
}

func TestMakeDaemonSetPositionDBAndSecrets(t *testing.T) {
	g := NewWithT(t)
	positionDB := corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/fluent-bit"}}
	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			DisableLogVolumes: true,
			PositionDB:        positionDB,
			Secrets:           []string{"es-credentials", "tls"},
		},
	}

	podSpec := MakeDaemonSet(fb, testLogPath).Spec.Template.Spec

	g.Expect(podSpec.Volumes[1:]).To(Equal([]corev1.Volume{
		{Name: "positions", VolumeSource: positionDB},
		secretVolume("es-credentials", "es-credentials"),
		secretVolume("tls", "tls"),
	}))
	g.Expect(podSpec.Containers[0].VolumeMounts[1:]).To(Equal([]corev1.VolumeMount{
		{Name: "positions", MountPath: "/fluent-bit/tail"},
		{Name: "es-credentials", ReadOnly: true, MountPath: "/fluent-bit/secrets/es-credentials"},
		{Name: "tls", ReadOnly: true, MountPath: "/fluent-bit/secrets/tls"},
	}))

	// Without a position DB no positions volume is added.
	fb.Spec.PositionDB = corev1.VolumeSource{}
	for _, v := range MakeDaemonSet(fb, testLogPath).Spec.Template.Spec.Volumes {
		g.Expect(v.Name).NotTo(Equal("positions"))
	}
}

func TestMakeDaemonSetLabels(t *testing.T) {
	g := NewWithT(t)
	fb := fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"from": "metadata"}},
	}

	// Without spec.labels the labels of the FluentBit resource are used.
	ds := MakeDaemonSet(fb, testLogPath)
	g.Expect(ds.Labels).To(Equal(map[string]string{"from": "metadata"}))
	g.Expect(ds.Spec.Selector.MatchLabels).To(Equal(ds.Labels))
	g.Expect(ds.Spec.Template.Labels).To(Equal(ds.Labels))

	// spec.labels replace them, they are not merged.
	fb.Spec.Labels = map[string]string{"from": "spec", "app": "fluent-bit"}
	ds = MakeDaemonSet(fb, testLogPath)
	g.Expect(ds.Labels).To(Equal(fb.Spec.Labels))
	g.Expect(ds.Spec.Selector.MatchLabels).To(Equal(fb.Spec.Labels))
	g.Expect(ds.Spec.Template.Labels).To(Equal(fb.Spec.Labels))
}

func TestMakeDaemonSetPortsAndEnv(t *testing.T) {
	g := NewWithT(t)
	fb := fluentbitv1alpha2.FluentBit{}

	container := MakeDaemonSet(fb, testLogPath).Spec.Template.Spec.Containers[0]
	g.Expect(container.Ports).To(Equal([]corev1.ContainerPort{{Name: "metrics", ContainerPort: 2020, Protocol: "TCP"}}))
	g.Expect(container.Env).To(HaveLen(2))
	g.Expect(container.Env[0].Name).To(Equal("NODE_NAME"))
	g.Expect(container.Env[0].ValueFrom.FieldRef.FieldPath).To(Equal("spec.nodeName"))
	g.Expect(container.Env[1].Name).To(Equal("HOST_IP"))
	g.Expect(container.Env[1].ValueFrom.FieldRef.FieldPath).To(Equal("status.hostIP"))

	fb.Spec.MetricsPort = 2021
	fb.Spec.Ports = []corev1.ContainerPort{{Name: "forward", ContainerPort: 24224, Protocol: "TCP"}}
	fb.Spec.EnvVars = []corev1.EnvVar{{Name: "CLUSTER", Value: "prod"}}

	container = MakeDaemonSet(fb, testLogPath).Spec.Template.Spec.Containers[0]
	g.Expect(container.Ports).To(Equal([]corev1.ContainerPort{
		{Name: "metrics", ContainerPort: 2021, Protocol: "TCP"},
		{Name: "forward", ContainerPort: 24224, Protocol: "TCP"},
	}))
	g.Expect(container.Env).To(HaveLen(3))
	g.Expect(container.Env[2]).To(Equal(corev1.EnvVar{Name: "CLUSTER", Value: "prod"}))
}

func TestMakeDaemonSetOptionalPodSettings(t *testing.T) {
	g := NewWithT(t)

	// Left unset, the optional settings keep their zero values so the API server defaults apply.
	ds := MakeDaemonSet(fluentbitv1alpha2.FluentBit{}, testLogPath)
	podSpec := ds.Spec.Template.Spec
	g.Expect(podSpec.RuntimeClassName).To(BeNil())
	g.Expect(podSpec.DNSPolicy).To(BeEmpty())
	g.Expect(podSpec.PriorityClassName).To(BeEmpty())
	g.Expect(podSpec.SchedulerName).To(BeEmpty())
	g.Expect(podSpec.Containers[0].Args).To(BeNil())
	g.Expect(podSpec.Containers[0].Command).To(BeNil())

	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			Image:             "fluent/fluent-bit:5.2.0",
			ImagePullPolicy:   corev1.PullIfNotPresent,
			RuntimeClassName:  "gvisor",
			DNSPolicy:         corev1.DNSClusterFirstWithHostNet,
			PriorityClassName: "system-node-critical",
			SchedulerName:     "custom-scheduler",
			Args:              []string{"--workdir=/fluent-bit/etc"},
			Command:           []string{"/fluent-bit/bin/fluent-bit"},
			HostNetwork:       true,
			NodeSelector:      map[string]string{"kubernetes.io/os": "linux"},
			Tolerations:       []corev1.Toleration{{Operator: corev1.TolerationOpExists}},
		},
	}
	podSpec = MakeDaemonSet(fb, testLogPath).Spec.Template.Spec
	g.Expect(podSpec.RuntimeClassName).To(HaveValue(Equal("gvisor")))
	g.Expect(podSpec.DNSPolicy).To(Equal(corev1.DNSClusterFirstWithHostNet))
	g.Expect(podSpec.PriorityClassName).To(Equal("system-node-critical"))
	g.Expect(podSpec.SchedulerName).To(Equal("custom-scheduler"))
	g.Expect(podSpec.HostNetwork).To(BeTrue())
	g.Expect(podSpec.NodeSelector).To(Equal(fb.Spec.NodeSelector))
	g.Expect(podSpec.Tolerations).To(Equal(fb.Spec.Tolerations))
	g.Expect(podSpec.Containers[0].Name).To(Equal("fluent-bit"))
	g.Expect(podSpec.Containers[0].Image).To(Equal("fluent/fluent-bit:5.2.0"))
	g.Expect(podSpec.Containers[0].ImagePullPolicy).To(Equal(corev1.PullIfNotPresent))
	g.Expect(podSpec.Containers[0].Args).To(Equal(fb.Spec.Args))
	g.Expect(podSpec.Containers[0].Command).To(Equal(fb.Spec.Command))
}
