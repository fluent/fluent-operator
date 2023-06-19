package operator

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
)

func MakeFluentdPVC(fd fluentdv1alpha1.Fluentd) *corev1.PersistentVolumeClaim {
	bufferStorage := fd.Spec.BufferVolume
	if bufferStorage == nil || bufferStorage.PersistentVolumeClaim == nil {
		return makeDefaultFluentdPVC(fd)
	}
	bufferPvc := bufferStorage.PersistentVolumeClaim.Spec

	labels := map[string]string{
		"app.kubernetes.io/name":      fd.Name,
		"app.kubernetes.io/instance":  "fluentd",
		"app.kubernetes.io/component": "fluentd",
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-buffer-pvc", fd.Name),
			Namespace: fd.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      bufferPvc.AccessModes,
			Resources:        bufferPvc.Resources,
			VolumeMode:       bufferPvc.VolumeMode,
			StorageClassName: bufferPvc.StorageClassName,
		},
	}
	return &pvc
}

func makeDefaultFluentdPVC(fd fluentdv1alpha1.Fluentd) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/name":      fd.Name,
		"app.kubernetes.io/instance":  "fluentd",
		"app.kubernetes.io/component": "fluentd",
	}

	r := corev1.ResourceRequirements{
		Requests: corev1.ResourceList(map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: resource.MustParse("1Gi")}),
	}

	fsmode := corev1.PersistentVolumeFilesystem

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-buffer-pvc", fd.Name),
			Namespace: fd.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:   r,
			VolumeMode:  &fsmode,
		},
	}
	return &pvc
}
