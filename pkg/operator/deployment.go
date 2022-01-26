package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "fluent.io/fluent-operator/apis/fluentd/v1alpha1"
)

const (
	SecretVolName    = "config"
	FluentdMountPath = "/fluentd/etc"
	BufferMountPath  = "/buffers"

	MetricsName = "metrics"

	MetricsPort int32 = 2021

	DefaultForwardPort int32 = 24424
	DefaultHttpPort    int32 = 9880

	DefaultForwardName = "forward"
	DefaultHttpName    = "http"

	InputForwardType = "forward"
	InputHttpType    = "http"
)

func MakeDeployment(fd fluentdv1alpha1.Fluentd) appsv1.Deployment {
	replicas := *fd.Spec.Replicas
	if replicas == 0 {
		replicas = 1
	}

	ports := makeDeploymentPorts(fd)

	labels := map[string]string{
		"app.kubernetes.io/name":      fd.Name,
		"app.kubernetes.io/instance":  "fluentd",
		"app.kubernetes.io/component": "fluentd",
	}

	if len(fd.Labels) > 0 {
		for k, v := range fd.Labels {
			if _, ok := labels[k]; !ok {
				labels[k] = v
			}
		}
	}

	dp := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fd.Name,
					Namespace: fd.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fd.Name,
					ImagePullSecrets:   fd.Spec.ImagePullSecrets,
					Volumes: []corev1.Volume{
						{
							Name: SecretVolName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fmt.Sprintf("%s-config", fd.Name),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "fluentd",
							Image:           fd.Spec.Image,
							Args:            fd.Spec.Args,
							ImagePullPolicy: fd.Spec.ImagePullPolicy,
							Ports:           ports,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      SecretVolName,
									ReadOnly:  true,
									MountPath: FluentdMountPath,
								},
							},
							Resources: fd.Spec.Resources,
							Env: []corev1.EnvVar{
								{
									Name:  "BUFFER_PATH",
									Value: BufferMountPath,
								},
							},
						},
					},
					NodeSelector: fd.Spec.NodeSelector,
					Tolerations:  fd.Spec.Tolerations,
					Affinity:     fd.Spec.Affinity,
				},
			},
		},
	}

	if fd.Spec.RuntimeClassName != "" {
		dp.Spec.Template.Spec.RuntimeClassName = &fd.Spec.RuntimeClassName
	}

	if fd.Spec.PriorityClassName != "" {
		dp.Spec.Template.Spec.PriorityClassName = fd.Spec.PriorityClassName
	}

	bufferVolName := fmt.Sprintf("%s-buffer", fd.Name)

	// Mount host or emptydir VolumeSource
	if fd.Spec.BufferVolume != nil && !fd.Spec.BufferVolume.DisableBufferVolume {
		bufferpvc := fd.Spec.BufferVolume

		if bufferpvc.HostPath != nil {
			dp.Spec.Template.Spec.Volumes = append(dp.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					HostPath: bufferpvc.HostPath,
				},
			})

			dp.Spec.Template.Spec.Containers[0].VolumeMounts = append(dp.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})
			return dp
		}

		if bufferpvc.EmptyDir != nil {
			dp.Spec.Template.Spec.Volumes = append(dp.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: bufferpvc.EmptyDir,
				},
			})

			dp.Spec.Template.Spec.Containers[0].VolumeMounts = append(dp.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})

			return dp
		}
	}

	// Bind pvc
	bufferPVCName := fmt.Sprintf("%s-buffer-pvc", fd.Name)
	pvcVol := corev1.Volume{
		Name: bufferPVCName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: bufferPVCName,
			},
		},
	}

	dp.Spec.Template.Spec.Volumes = append(dp.Spec.Template.Spec.Volumes, pvcVol)

	dp.Spec.Template.Spec.Containers[0].VolumeMounts = append(dp.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      bufferPVCName,
		MountPath: BufferMountPath,
	})

	return dp

}

func makeDeploymentPorts(fd fluentdv1alpha1.Fluentd) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          MetricsName,
			ContainerPort: MetricsPort,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	// read inputs definition from globalInputs
	globalInputs := fd.Spec.GlobalInputs
	for _, input := range globalInputs {
		if *input.InputCommon.Type == InputForwardType && input.Forward != nil {
			forwardPort := *input.Forward.Port
			if forwardPort == 0 {
				forwardPort = DefaultForwardPort
			}
			ports = append(ports, corev1.ContainerPort{
				Name:          DefaultForwardName,
				ContainerPort: forwardPort,
				Protocol:      corev1.ProtocolTCP,
			})
			continue
		}
		if *input.InputCommon.Type == InputHttpType && input.Http != nil {
			httpPort := *input.Forward.Port
			if httpPort == 0 {
				httpPort = DefaultHttpPort
			}
			ports = append(ports, corev1.ContainerPort{
				Name:          DefaultHttpName,
				ContainerPort: httpPort,
				Protocol:      corev1.ProtocolTCP,
			})
		}
	}

	return ports
}
