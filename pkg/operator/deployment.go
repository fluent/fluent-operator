package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2"
)

func MakeDeployment(co fluentbitv1alpha2.Collector) appsv1.Deployment {
	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      co.Name,
			Namespace: co.Namespace,
			Labels:    co.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: co.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        co.Name,
					Namespace:   co.Namespace,
					Labels:      co.Labels,
					Annotations: co.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: co.Name,
					ImagePullSecrets:   co.Spec.ImagePullSecrets,
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: co.Spec.FluentBitConfigName,
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
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           co.Spec.Image,
							Args:            co.Spec.Args,
							ImagePullPolicy: co.Spec.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
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
							},
							VolumeMounts: []corev1.VolumeMount{
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
							Resources: co.Spec.Resources,
						},
					},
					NodeSelector:    co.Spec.NodeSelector,
					Tolerations:     co.Spec.Tolerations,
					Affinity:        co.Spec.Affinity,
					SecurityContext: co.Spec.SecurityContext,
					HostNetwork:     co.Spec.HostNetwork,
				},
			},
		},
	}

	if co.Spec.RuntimeClassName != "" {
		deploy.Spec.Template.Spec.RuntimeClassName = &co.Spec.RuntimeClassName
	}

	if co.Spec.PriorityClassName != "" {
		deploy.Spec.Template.Spec.PriorityClassName = co.Spec.PriorityClassName
	}

	if co.Spec.Volumes != nil {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, co.Spec.Volumes...)
	}
	if co.Spec.VolumesMounts != nil {
		deploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(deploy.Spec.Template.Spec.Containers[0].VolumeMounts, co.Spec.VolumesMounts...)
	}

	// Mount Position DB
	if co.Spec.PositionDB != (corev1.VolumeSource{}) {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         "positions",
			VolumeSource: co.Spec.PositionDB,
		})
		deploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(deploy.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "positions",
			MountPath: "/fluent-bit/tail",
		})
	}

	// Mount Secrets
	for _, secret := range co.Spec.Secrets {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: secret,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret,
				},
			},
		})
		deploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(deploy.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      secret,
			ReadOnly:  true,
			MountPath: fmt.Sprintf("/fluent-bit/secrets/%s", secret),
		})
	}

	return deploy
}
