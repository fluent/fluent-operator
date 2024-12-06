package operator

import (
	"testing"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var (
	logPath string = "/var/lib/docker/containers"
)

func TestDaemonsetCustomVolumes(t *testing.T) {
	g := NewWithT(t)

	customVolumes := []corev1.Volume{
		{
			Name: "pvc-volume",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "pvc",
				},
			},
		},
		{
			Name: "emptydir-volume",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "nfs-volume",
			VolumeSource: corev1.VolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server:   "nfs-server",
					Path:     "/mnt",
					ReadOnly: true,
				},
			},
		},
	}
	customVolumeMounts := []corev1.VolumeMount{
		{
			Name:      "pvc-volume",
			MountPath: "/mnt/pvc",
		},
		{
			Name:      "emptydir-volume",
			MountPath: "/mnt/emptydir",
		},
		{
			Name:      "nfs-volume",
			MountPath: "/mnt/nfs",
		},
	}

	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			DisableLogVolumes:   true,
			Volumes:             customVolumes,
			VolumesMounts:       customVolumeMounts,
			FluentBitConfigName: "fb-config",
		},
	}

	fbWithLogDisabled := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			DisableLogVolumes: false,
			VolumesMounts:     customVolumeMounts,
			Volumes:           customVolumes,
		},
	}

	ds := MakeDaemonSet(fb, logPath)
	dsWithLogDisabled := MakeDaemonSet(fbWithLogDisabled, logPath)
	g.Expect(ds.Spec.Template.Spec.Volumes).Should(ContainElements(customVolumes))
	g.Expect(ds.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElements(customVolumeMounts))
	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Volumes).Should(ContainElements(customVolumes))
	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElements(customVolumeMounts))
}

func TestDaemonsetDefaultVolumes(t *testing.T) {
	g := NewWithT(t)

	internalMountPropagation := corev1.MountPropagationNone
	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			DisableLogVolumes:   false,
			FluentBitConfigName: "fb-config",
		},
	}
	fbVolumesLogDisabled := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			DisableLogVolumes:   true,
			FluentBitConfigName: "fb-config",
		},
	}

	expectedVolumesLogDisabled := []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: fb.Spec.FluentBitConfigName,
				},
			},
		},
	}
	expectedVolumeMountsLogDisabled := []corev1.VolumeMount{
		{
			Name:      "config",
			ReadOnly:  true,
			MountPath: "/fluent-bit/config",
		},
	}

	expectedVolumes := []corev1.Volume{
		{
			Name: "varlibcontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: logPath,
				},
			},
		},
		{
			Name: "varlogs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log",
				},
			},
		},
		{
			Name: "systemd",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log/journal",
				},
			},
		},
	}
	expectedVolumeMounts := []corev1.VolumeMount{
		{
			Name:             "varlibcontainers",
			ReadOnly:         true,
			MountPath:        logPath,
			MountPropagation: &internalMountPropagation,
		},

		{
			Name:             "varlogs",
			ReadOnly:         true,
			MountPath:        "/var/log/",
			MountPropagation: &internalMountPropagation,
		},
		{
			Name:             "systemd",
			ReadOnly:         true,
			MountPath:        "/var/log/journal",
			MountPropagation: &internalMountPropagation,
		},
	}

	ds := MakeDaemonSet(fb, logPath)
	dsWithLogDisabled := MakeDaemonSet(fbVolumesLogDisabled, logPath)

	g.Expect(ds.Spec.Template.Spec.Volumes).Should(ContainElements(expectedVolumes))
	g.Expect(ds.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElements(expectedVolumeMounts))

	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Volumes).Should(ContainElements(expectedVolumesLogDisabled))
	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Volumes).ShouldNot(ContainElements(expectedVolumes))
	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Containers[0].VolumeMounts).Should(ContainElements(expectedVolumeMountsLogDisabled))
	g.Expect(dsWithLogDisabled.Spec.Template.Spec.Containers[0].VolumeMounts).ShouldNot(ContainElements(expectedVolumeMounts))
}

func TestDaemonsetLabelsAndAnnotations(t *testing.T) {
	g := NewWithT(t)

	annotations := map[string]string{
		"test_a": "value_a",
		"test_b": "value_b",
	}
	labels := map[string]string{
		"test_a": "value_a",
		"test_b": "value_b",
	}
	fb := fluentbitv1alpha2.FluentBit{
		Spec: fluentbitv1alpha2.FluentBitSpec{
			Annotations: annotations,
			Labels:      labels,
		},
	}

	ds := MakeDaemonSet(fb, logPath)

	g.Expect(ds.Spec.Template.Annotations).Should(Equal(annotations))
	g.Expect(ds.Labels).Should(Equal(labels))
	g.Expect(ds.Spec.Template.Labels).Should(Equal(labels))
	g.Expect(ds.Spec.Selector.MatchLabels).Should(Equal(labels))
}
