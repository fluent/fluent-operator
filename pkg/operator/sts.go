package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
)

const (
	SecretVolName    = "config"
	FluentdMountPath = "/fluentd/etc"
	BufferMountPath  = "/buffers"

	MetricsName = "metrics"

	MetricsPort int32 = 2021

	DefaultForwardPort int32 = 24424
	DefaultHttpPort    int32 = 9880
	// 101 is the fsGroup that fluentd runs as in the kubesphere image
	DefaultFsGroup int64 = 101

	DefaultForwardName = "forward"
	DefaultHttpName    = "http"

	InputForwardType = "forward"
	InputHttpType    = "http"
)

func MakeStatefulSet(fd fluentdv1alpha1.Fluentd) *appsv1.StatefulSet {
	replicas := *fd.Spec.Replicas
	if replicas == 0 {
		replicas = 1
	}

	ports := makeFluentdPorts(fd)

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

	defaultFsGroup := DefaultFsGroup

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fd.Name,
			Namespace:   fd.Namespace,
			Labels:      labels,
			Annotations: fd.Annotations,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fd.Name,
					Namespace:   fd.Namespace,
					Labels:      labels,
					Annotations: fd.Spec.Annotations,
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
							ReadinessProbe: fd.Spec.ReadinessProbe,
							LivenessProbe:  fd.Spec.LivenessProbe,
						},
					},
					NodeSelector: fd.Spec.NodeSelector,
					Tolerations:  fd.Spec.Tolerations,
					Affinity:     fd.Spec.Affinity,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &defaultFsGroup,
					},
				},
			},
		},
	}

	if fd.Spec.RuntimeClassName != "" {
		sts.Spec.Template.Spec.RuntimeClassName = &fd.Spec.RuntimeClassName
	}

	if fd.Spec.PriorityClassName != "" {
		sts.Spec.Template.Spec.PriorityClassName = fd.Spec.PriorityClassName
	}

	if fd.Spec.Volumes != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, fd.Spec.Volumes...)
	}

	if fd.Spec.VolumeMounts != nil {
		sts.Spec.Template.Spec.Containers[0].VolumeMounts = append(sts.Spec.Template.Spec.Containers[0].VolumeMounts, fd.Spec.VolumeMounts...)
	}

	if fd.Spec.VolumeClaimTemplates != nil {
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, fd.Spec.VolumeClaimTemplates...)
	}

	if fd.Spec.EnvVars != nil {
		sts.Spec.Template.Spec.Containers[0].Env = append(sts.Spec.Template.Spec.Containers[0].Env, fd.Spec.EnvVars...)
	}

	if fd.Spec.SecurityContext != nil {
		sts.Spec.Template.Spec.SecurityContext = fd.Spec.SecurityContext
	}

	if fd.Spec.SchedulerName != "" {
		sts.Spec.Template.Spec.SchedulerName = fd.Spec.SchedulerName
	}

	// Mount host or emptydir VolumeSource
	if fd.Spec.BufferVolume != nil && !fd.Spec.BufferVolume.DisableBufferVolume {
		bufferVolName := fmt.Sprintf("%s-buffer", fd.Name)
		bufferpv := fd.Spec.BufferVolume

		if bufferpv.HostPath != nil {
			sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					HostPath: bufferpv.HostPath,
				},
			})

			sts.Spec.Template.Spec.Containers[0].VolumeMounts = append(sts.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})
			return &sts
		}

		if bufferpv.EmptyDir != nil {
			sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: bufferpv.EmptyDir,
				},
			})

			sts.Spec.Template.Spec.Containers[0].VolumeMounts = append(sts.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})

			return &sts
		}
	}

	if fd.Spec.BufferVolume == nil || !fd.Spec.BufferVolume.DisableBufferVolume {
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *MakeFluentdPVC(fd))

		sts.Spec.Template.Spec.Containers[0].VolumeMounts = append(sts.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("%s-buffer-pvc", fd.Name),
			MountPath: BufferMountPath,
		})
	}
	return &sts
}

func makeFluentdPorts(fd fluentdv1alpha1.Fluentd) []corev1.ContainerPort {
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
		if input.Forward != nil {
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
		if input.Http != nil {
			httpPort := *input.Http.Port
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
