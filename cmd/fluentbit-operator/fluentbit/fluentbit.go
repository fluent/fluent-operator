package fluentbit

import (
	"bytes"
	"kubesphere.io/fluentbit-operator/cmd/fluentbit-operator/sdkdecorator"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	extensionv1 "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sync"
	"text/template"
)

// OwnerDeployment of the daemonset
var OwnerDeployment metav1.Object
var config *fluentBitDeploymentConfig

// ConfigLock used for AppConfig
var ConfigLock sync.Mutex

func initConfig(labels map[string]string) *fluentBitDeploymentConfig {
	if config == nil {
		config = &fluentBitDeploymentConfig{
			Name:      "fluent-bit",
			Namespace: viper.GetString("fluent-bit.namespace"),
			Labels:    labels,
		}
		config.Labels["app"] = "fluent-bit"
	}
	return config
}

// InitFluentBit initialize fluent-bit
func InitFluentBit(labels map[string]string) {
	cfg := initConfig(labels)
	if !checkIfDeamonSetExist(cfg) {
		logrus.Info("Deploying fluent-bit")
		if viper.GetBool("fluentbit-operator.rbac") {
			sa := newServiceAccount(cfg)
			err := controllerutil.SetControllerReference(OwnerDeployment, sa, scheme.Scheme)
			logrus.Error(err)
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(sa)
			cr := newClusterRole(cfg)
			err = controllerutil.SetControllerReference(OwnerDeployment, cr, scheme.Scheme)
			logrus.Error(err)
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(cr)
			crb := newClusterRoleBinding(cfg)
			err = controllerutil.SetControllerReference(OwnerDeployment, crb, scheme.Scheme)
			logrus.Error(err)
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(crb)
		}
		cfgMap, err := newFluentBitConfig(cfg)
		if err != nil {
			logrus.Error(err)
		}
		err = controllerutil.SetControllerReference(OwnerDeployment, cfgMap, scheme.Scheme)
		logrus.Error(err)
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(cfgMap)
		CreateOrUpdateAppConfig("", "", "")
		ds := newFluentBitDaemonSet(cfg)
		err = controllerutil.SetControllerReference(OwnerDeployment, ds, scheme.Scheme)
		logrus.Error(err)
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(ds)
		logrus.Info("Fluent-bit deployed successfully")
	}
}

// DeleteFluentBit deletes fluent-bit if it exists
func DeleteFluentBit(labels map[string]string) {
	cfg := initConfig(labels)
	if checkIfDeamonSetExist(cfg) {
		logrus.Info("Deleting fluent-bit")
		if viper.GetBool("fluentbit-operator.rbac") {
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newServiceAccount(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newClusterRole(cfg))
			sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newClusterRoleBinding(cfg))
		}
		cfgMap, err := newFluentBitConfig(cfg)
		if err != nil {
			logrus.Error(err)
		}
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(cfgMap)
		DeleteAppConfig()
		foregroundDeletion := metav1.DeletePropagationForeground
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(newFluentBitDaemonSet(cfg),
			sdk.WithDeleteOptions(&metav1.DeleteOptions{
				PropagationPolicy: &foregroundDeletion,
			}))
		logrus.Info("Fluent-bit deleted successfully")
	}
}

type fluentBitDeploymentConfig struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type fluentBitConfig struct {
	Namespace string
}

func newServiceAccount(cr *fluentBitDeploymentConfig) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentbit",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
	}
}

func newClusterRole(cr *fluentBitDeploymentConfig) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "LoggingRole",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs: []string{
					"get",
				},
				APIGroups: []string{""},
				Resources: []string{
					"pods",
				},
			},
		},
	}
}

func newClusterRoleBinding(cr *fluentBitDeploymentConfig) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentbit",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "fluentbit",
				Namespace: cr.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "LoggingRole",
		},
	}
}

func generateConfig(input fluentBitConfig) (*string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentBitConfigTemplate)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return nil, err
	}
	outputString := output.String()
	return &outputString, nil
}

func generateSettings(input fluentBitConfig) (*string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentBitSettingsTemplate)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return nil, err
	}
	outputString := output.String()
	return &outputString, nil
}

// DeleteAppConfig thread safe config management
func DeleteAppConfig() {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit-app-config",
			Namespace: config.Namespace,
			Labels:    config.Labels,
		},
	}
	ConfigLock.Lock()
	defer ConfigLock.Unlock()
	sdkdecorator.CallSdkFunctionWithLogging(sdk.Delete)(configMap)
}

// CreateOrUpdateAppConfig idempotent thread safe config management
func CreateOrUpdateAppConfig(name string, appConfig string, appSettings string) {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit-app-config",
			Namespace: config.Namespace,
			Labels:    config.Labels,
		},
	}
	// Lock for shared fluent-bit config resource
	ConfigLock.Lock()
	defer ConfigLock.Unlock()
	err := sdk.Get(configMap)
	if err != nil && !apierrors.IsNotFound(err) {
		// Something unexpected happened
		logrus.Error(err)
		return
	}
	// Do the changes
	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}
	if name != "" && appConfig != "" {
		configMap.Data[name+".conf"] = appConfig
	}
	if name != "" && appSettings != "" {
		configMap.Data["settings.json"] = appSettings
	}
	// The resource not Found so we create it
	if err != nil {
		err = controllerutil.SetControllerReference(OwnerDeployment, configMap, scheme.Scheme)
		logrus.Error(err)
		sdkdecorator.CallSdkFunctionWithLogging(sdk.Create)(configMap)
		return
	}
	// No error we go for update
	sdkdecorator.CallSdkFunctionWithLogging(sdk.Update)(configMap)
}

func newFluentBitConfig(cr *fluentBitDeploymentConfig) (*corev1.ConfigMap, error) {
	input := fluentBitConfig{
		Namespace: cr.Namespace,
	}
	config, err := generateConfig(input)
	if err != nil {
		return nil, err
	}
	settings, err := generateSettings(input)
	if err != nil {
		return nil, err
	}
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit-config",
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},

		Data: map[string]string{
			"fluent-bit.conf": *config,
			"settings.json": *settings,
		},
	}
	return configMap, nil
}

func checkIfDeamonSetExist(cr *fluentBitDeploymentConfig) bool {
	fluentbitDaemonSet := &extensionv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Labels:    cr.Labels,
			Namespace: cr.Namespace,
		},
	}
	if err := sdk.Get(fluentbitDaemonSet); err != nil {
		logrus.Info("FluentBit DaemonSet does not exists!")
		logrus.Error(err)
		return false
	}
	logrus.Info("FluentBit DaemonSet already exists!")
	return true
}

func newConfigMapReloader() *corev1.Container {
	return &corev1.Container{
		Name:  "config-reloader",
		Image: viper.GetString("configmap-reload.image"),
		Args: []string{
			"-volume-dir=/fluent-bit/app-config/",
			"-webhook-url=http://127.0.0.1:24444/api/config.reload",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "app-config",
				MountPath: "/fluent-bit/app-config/",
			},
		},
	}
}

func generateVolumeMounts() (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
		{
			Name:      "varlibcontainers",
			ReadOnly:  true,
			MountPath: viper.GetString("fluent-bit.containersLogMountedPath"),
		},
		{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   "fluent-bit.conf",
		},
		{
			Name:      "app-config",
			MountPath: "/fluent-bit/app-config/",
		},
		{
			Name:      "positions",
			MountPath: "/tail-db",
		},
		{
			Name:      "varlogs",
			ReadOnly:  true,
			MountPath: "/var/log/",
		},
	}
	return
}

func generateVolume() (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "varlibcontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: viper.GetString("fluent-bit.containersLogMountedPath"),
				},
			},
		},
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluent-bit-config",
					},
				},
			},
		},
		{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluent-bit-app-config",
					},
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
			Name: "positions",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
	return
}

// TODO in case of rbac add created serviceAccount name
func newFluentBitDaemonSet(cr *fluentBitDeploymentConfig) *extensionv1.DaemonSet {
	return &extensionv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: extensionv1.DaemonSetSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.Labels,
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/api/v1/metrics/prometheus",
						"prometheus.io/port":   "2020",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "fluentbit",
					Volumes:            generateVolume(),
					Containers: []corev1.Container{
						{
							// TODO move to configuration
							Name:  "fluent-bit",
							Image: viper.GetString("fluent-bit.image"),
							// TODO get from config translate to const
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
							// TODO Get this from config
							Resources: corev1.ResourceRequirements{
								Limits:   nil,
								Requests: nil,
							},
							VolumeMounts: generateVolumeMounts(),
						},
						*newConfigMapReloader(),
					},
					Tolerations: []corev1.Toleration{
						{
							Key: "CriticalAddonsOnly",
							Operator: "Exists",
						},
						{
							Key: "node-role.kubernetes.io/master",
							Effect: "NoSchedule",
						},
						{
							Key: "dedicated",
							Operator: "Exists",
						},
						{
							Key: "node.cloudprovider.kubernetes.io/uninitialized",
							Value: "true",
							Effect: "NoSchedule",
						},
					},
				},
			},
		},
	}
}
