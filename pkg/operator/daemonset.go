package operator

import (
	"fmt"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeDaemonSet(fb fluentbitv1alpha2.FluentBit, logPath string) appsv1.DaemonSet {
	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fb.Name,
			Namespace:   fb.Namespace,
			Labels:      fb.Labels,
			Annotations: fb.Annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: fb.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fb.Name,
					Namespace:   fb.Namespace,
					Labels:      fb.Labels,
					Annotations: fb.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fb.Name,
					ImagePullSecrets:   fb.Spec.ImagePullSecrets,
					Volumes: []corev1.Volume{
						{
							Name: "varlibcontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: logPath,
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fb.Spec.FluentBitConfigName,
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
					},
					InitContainers: fb.Spec.InitContainers,
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           fb.Spec.Image,
							ImagePullPolicy: fb.Spec.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
							ReadinessProbe: fb.Spec.ReadinessProbe,
							LivenessProbe:  fb.Spec.LivenessProbe,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "varlibcontainers",
									ReadOnly:  true,
									MountPath: logPath,
								},
								{
									Name:      "config",
									ReadOnly:  true,
									MountPath: "/fluent-bit/config",
								},
								{
									Name:      "varlogs",
									ReadOnly:  true,
									MountPath: "/var/log/",
								},
								{
									Name:      "systemd",
									ReadOnly:  true,
									MountPath: "/var/log/journal",
								},
							},
							Resources: fb.Spec.Resources,
						},
					},
					NodeSelector:    fb.Spec.NodeSelector,
					Tolerations:     fb.Spec.Tolerations,
					Affinity:        fb.Spec.Affinity,
					SecurityContext: fb.Spec.SecurityContext,
					HostNetwork:     fb.Spec.HostNetwork,
				},
			},
		},
	}

	if fb.Spec.Args != nil {
		ds.Spec.Template.Spec.Containers[0].Args = fb.Spec.Args
	}

	if fb.Spec.Command != nil {
		ds.Spec.Template.Spec.Containers[0].Command = fb.Spec.Command
	}

	if fb.Spec.Ports != nil {
		ds.Spec.Template.Spec.Containers[0].Ports = append(ds.Spec.Template.Spec.Containers[0].Ports, fb.Spec.Ports...)
	}

	if fb.Spec.EnvVars != nil {
		ds.Spec.Template.Spec.Containers[0].Env = append(ds.Spec.Template.Spec.Containers[0].Env, fb.Spec.EnvVars...)
	}

	if fb.Spec.RuntimeClassName != "" {
		ds.Spec.Template.Spec.RuntimeClassName = &fb.Spec.RuntimeClassName
	}

	if fb.Spec.PriorityClassName != "" {
		ds.Spec.Template.Spec.PriorityClassName = fb.Spec.PriorityClassName
	}

	if fb.Spec.Volumes != nil {
		ds.Spec.Template.Spec.Volumes = append(ds.Spec.Template.Spec.Volumes, fb.Spec.Volumes...)
	}
	if fb.Spec.VolumesMounts != nil {
		ds.Spec.Template.Spec.Containers[0].VolumeMounts = append(ds.Spec.Template.Spec.Containers[0].VolumeMounts, fb.Spec.VolumesMounts...)
	}

	// Mount Position DB
	if fb.Spec.PositionDB != (corev1.VolumeSource{}) {
		ds.Spec.Template.Spec.Volumes = append(ds.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         "positions",
			VolumeSource: fb.Spec.PositionDB,
		})
		ds.Spec.Template.Spec.Containers[0].VolumeMounts = append(ds.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "positions",
			MountPath: "/fluent-bit/tail",
		})
	}

	// Mount Secrets
	for _, secret := range fb.Spec.Secrets {
		ds.Spec.Template.Spec.Volumes = append(ds.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: secret,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret,
				},
			},
		})
		ds.Spec.Template.Spec.Containers[0].VolumeMounts = append(ds.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      secret,
			ReadOnly:  true,
			MountPath: fmt.Sprintf("/fluent-bit/secrets/%s", secret),
		})
	}

	return ds
}
