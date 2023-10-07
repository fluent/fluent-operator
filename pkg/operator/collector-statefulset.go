package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
)

var (
	// DefaultBufferPath defines the buffer path for single process
	DefaultBufferPath = "/buffers/fluentbit/log"
)

func MakefbStatefulset(co fluentbitv1alpha2.Collector) *appsv1.StatefulSet {
	statefulset := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      co.Name,
			Namespace: co.Namespace,
			Labels:    co.Labels,
		},
		Spec: appsv1.StatefulSetSpec{
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
		statefulset.Spec.Template.Spec.RuntimeClassName = &co.Spec.RuntimeClassName
	}

	if co.Spec.PriorityClassName != "" {
		statefulset.Spec.Template.Spec.PriorityClassName = co.Spec.PriorityClassName
	}

	if co.Spec.SchedulerName != "" {
		statefulset.Spec.Template.Spec.SchedulerName = co.Spec.SchedulerName
	}

	if co.Spec.Volumes != nil {
		statefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, co.Spec.Volumes...)
	}
	if co.Spec.VolumesMounts != nil {
		statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, co.Spec.VolumesMounts...)
	}

	if co.Spec.Ports != nil {
		statefulset.Spec.Template.Spec.Containers[0].Ports = append(statefulset.Spec.Template.Spec.Containers[0].Ports, co.Spec.Ports...)
	}

	// Mount Secrets
	for _, secret := range co.Spec.Secrets {
		statefulset.Spec.Template.Spec.Volumes = append(statefulset.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: secret,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret,
				},
			},
		})
		statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      secret,
			ReadOnly:  true,
			MountPath: fmt.Sprintf("/fluent-bit/secrets/%s", secret),
		})
	}

	// Bind pvc
	statefulset.Spec.VolumeClaimTemplates = append(statefulset.Spec.VolumeClaimTemplates, MakeFluentbitPVC(co))
	statefulset.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulset.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      fmt.Sprintf("%s-buffer-pvc", co.Name),
		MountPath: FluentbitBufferMountPath(co),
	})

	return &statefulset
}

func MakeFluentbitPVC(co fluentbitv1alpha2.Collector) corev1.PersistentVolumeClaim {
	bufferStorage := co.Spec.PersistentVolumeClaim
	if bufferStorage == nil {
		return makeDefaultFluentbitPVC(co)
	}
	bufferPvc := bufferStorage.Spec

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-buffer-pvc", co.Name),
			Namespace: co.Namespace,
			Labels:    co.Labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: bufferPvc.AccessModes,
			Resources:   bufferPvc.Resources,
			VolumeMode:  bufferPvc.VolumeMode,
		},
	}
	return pvc
}

func makeDefaultFluentbitPVC(co fluentbitv1alpha2.Collector) corev1.PersistentVolumeClaim {

	r := corev1.ResourceRequirements{
		Requests: corev1.ResourceList(map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resource.MustParse("1Gi")}),
	}

	fsmode := corev1.PersistentVolumeFilesystem

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-buffer-pvc", co.Name),
			Namespace: co.Namespace,
			Labels:    co.Labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:   r,
			VolumeMode:  &fsmode,
		},
	}
	return pvc
}

func FluentbitBufferMountPath(co fluentbitv1alpha2.Collector) string {
	bufferPath := co.Spec.BufferPath
	if bufferPath != nil {
		return *bufferPath
	} else {
		return DefaultBufferPath
	}
}
