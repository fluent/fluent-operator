package stub

import (
	"context"
	"kubesphere.io/fluentbit-operator/cmd/fluentbit-operator/fluentbit"
	"kubesphere.io/fluentbit-operator/pkg/apis/fluentbit/v1alpha1"
	"kubesphere.io/fluentbit-operator/pkg/plugins"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler creates a new Handler struct
func NewHandler(namepsace string) sdk.Handler {
	return &Handler{
		NameSpace: namepsace,
	}
}

// Handler struct
type Handler struct {
	NameSpace string
}

// Handle every event set up by the watcher
func (h *Handler) Handle(ctx context.Context, event sdk.Event) (err error) {
	switch o := event.Object.(type) {
	case *v1alpha1.FluentBitOperator:
		if event.Deleted {
			logrus.Infof("Delete CRD: %s", o.Name)
			deleteFromConfigMap(o.Name)
			return
		}
		logrus.Infof("New CRD arrived %#v", o)
		logrus.Info("Generating configuration.")
		name, config, settings := generateFluentbitConfigAndSettings(o, h.NameSpace)
		if config != "" && name != "" {
			fluentbit.CreateOrUpdateAppConfig(name, config, settings)
		}
	}
	return
}

func deleteFromConfigMap(name string) {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent-bit-app-config",
			Namespace: "default",
		},
	}
	err := sdk.Get(configMap)
	if err != nil {
		logrus.Error(err)
	}
	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}
	delete(configMap.Data, name+".conf")
	delete(configMap.Data, "settings.json")
	err = sdk.Update(configMap)
	if err != nil {
		logrus.Error(err)
	}
}

//
func generateFluentbitConfigAndSettings(crd *v1alpha1.FluentBitOperator, namespace string) (string, string, string) {
	var finalConfig string
	// Generate service
	for _, service := range crd.Spec.Service {
		logrus.Info("Applying service")
		values, err := plugins.GetDefaultValues(service.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		config, err := v1alpha1.RenderPlugin(service, values, namespace, "[SERVICE]")
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		finalConfig += config
	}

	// Generate input
	for _, input := range crd.Spec.Input {
		logrus.Info("Applying input")
		values, err := plugins.GetDefaultValues(input.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		config, err := v1alpha1.RenderPlugin(input, values, namespace, "[INPUT]")
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		finalConfig += config
	}

	// Generate filters
	for _, filter := range crd.Spec.Filter {
		logrus.Info("Applying filter")
		values, err := plugins.GetDefaultValues(filter.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		config, err := v1alpha1.RenderPlugin(filter, values, namespace, "[FILTER]")
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		finalConfig += config
	}

	// Generate output
	for _, output := range crd.Spec.Output {
		logrus.Info("Applying output")
		values, err := plugins.GetDefaultValues(output.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		config, err := v1alpha1.RenderPlugin(output, values, namespace, "[OUTPUT]")
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		finalConfig += config
	}

	var finalSettings string
	// Generate settings
	settings := map[string]string{}
	for _, setting := range crd.Spec.Settings {
		logrus.Info("Applying settings")
		values, err := plugins.GetDefaultValues(setting.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", "", ""
		}
		v1alpha1.ProcessSettings(setting, settings, values, namespace)
	}
	finalSettings, err := v1alpha1.RenderSettings(settings)
	if err != nil {
		logrus.Infof("Error in rendering template: %s", err)
		return "", "", ""
	}

	return crd.Name, finalConfig, finalSettings
}
