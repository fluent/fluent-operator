package main

import (
	"context"
	"runtime"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"kubesphere.io/fluentbit-operator/pkg/stub"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"kubesphere.io/fluentbit-operator/cmd/fluentbit-operator/fluentbit"
	"os"
)

func printVersion(namespace string) {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
	logrus.Infof("Operator namespace: %s", namespace)
}

func getConfiguration() *operatorConfig {
	return &operatorConfig{}
}

// GlobalLabels to use with generated deployments
var GlobalLabels = map[string]string{
	"chart":   "",
	"release": "",
}

func main() {

	const (
		operatorNamespace   = "WATCH_NAMESPACE"
		operatorResource    = "logging.kubesphere.io/v1alpha1"
		kind                = "FluentBit"
		kubernetesPodName   = "KUBERNETES_POD_NAME"
		kubernetesNamespace = "KUBERNETES_NAMESPACE"
	)
	podNamespace := os.Getenv(kubernetesNamespace)
	podName := os.Getenv(kubernetesPodName)
	logrus.Infof("Gettint current environment: ns: %q pod: %q", podNamespace, podName)
	pod, err := GetSelf(podName, podNamespace)
	if err != nil {
		logrus.Error(err.Error())
	}
	obj, err := GetDeployment(pod, pod.Namespace)
	if err != nil {
		logrus.Error(err.Error())
	}
	deploymentLabels := obj.GetLabels()
	GlobalLabels["chart"] = deploymentLabels["chart"]
	GlobalLabels["release"] = deploymentLabels["release"]
	fluentbit.OwnerDeployment = obj
	ns := os.Getenv(operatorNamespace)
	printVersion(ns)
	resyncPeriod := 0
	logrus.Infof("Watching %s, %s, %s, %d", operatorResource, kind, ns, resyncPeriod)
	sdk.Watch(operatorResource, kind, ns, resyncPeriod)
	err = viper.BindEnv("kubernetesNamespace", "KUBERNETES_NAMESPACE")
	if err != nil {
		logrus.Error(err)
	}
	// Init resources
	Init()

	sdk.Handle(stub.NewHandler(viper.GetString("kubernetesNamespace")))
	sdk.Run(context.TODO())
}

type operatorConfig struct {
}
