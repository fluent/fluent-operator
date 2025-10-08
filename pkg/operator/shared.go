package operator

import (
	"fmt"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// configureBufferVolume configures the buffer volume for fluentd.
// Returns true if the caller should return early.
func configureBufferVolume(fd fluentdv1alpha1.Fluentd, specTemplateSpec *corev1.PodSpec) bool {
	if fd.Spec.BufferVolume != nil && !fd.Spec.BufferVolume.DisableBufferVolume {
		bufferVolName := fmt.Sprintf("%s-buffer", fd.Name)
		bufferpv := fd.Spec.BufferVolume

		if bufferpv.HostPath != nil {
			specTemplateSpec.Volumes = append(specTemplateSpec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					HostPath: bufferpv.HostPath,
				},
			})

			specTemplateSpec.Containers[0].VolumeMounts = append(specTemplateSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})

			return true
		}

		if bufferpv.EmptyDir != nil {
			specTemplateSpec.Volumes = append(specTemplateSpec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: bufferpv.EmptyDir,
				},
			})

			specTemplateSpec.Containers[0].VolumeMounts = append(specTemplateSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})

			return true
		}
	}
	return false
}
