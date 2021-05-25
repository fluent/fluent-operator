package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/output"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	Namespace = "kubesphere-logging-system"
	Name      = "fluent-bit-output-config"
	Key       = "outputs"

	Elasticsearch = "fluentbit-output-es"
	Kafka         = "fluentbit-output-kafka"
	Fluentd       = "fluentbit-output-forward"
)

type oldConfig struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Params []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameters"`
	ID         string `json:"id"`
	Enable     bool   `json:"enable"`
	Updatetime string `json:"updatetime"`
}

var (
	masterURL      string
	kubeconfigPath string

	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = logging.AddToScheme(scheme)
}

func main() {
	flag.StringVar(&kubeconfigPath, "kubeconfigPath", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "masterURL", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()

	// Create k8s client
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		klog.Errorf("Error building kubeconfig: %s", err.Error())
		return
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Error building kubernetes clientset: %s", err.Error())
		return
	}

	// Read old config from the ConfigMap `fluent-bit-output-config`.
	cm, err := kubeClient.CoreV1().ConfigMaps(Namespace).Get(Name, corev1.GetOptions{})
	if err != nil {
		klog.Errorf("Failed to find configmap: %s/%s", Namespace, Name)
		return
	}

	// Decode config data
	src := make([]oldConfig, 0)
	err = jsoniter.UnmarshalFromString(cm.Data[Key], &src)
	if err != nil {
		klog.Error("Failed to decode config data")
		return
	}

	// Create client for fluent bit operator CRDs
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme})
	if err != nil {
		klog.Errorf("Error building fluent bit operator clientset: %s", err.Error())
		return
	}
	fboClient := mgr.GetClient()

	// Migrating
	for _, item := range src {
		var out logging.Output

		switch item.Name {
		case Elasticsearch:
			out = makeElasticsearchOutput(item)
		case Kafka:
			out = makeKafkaOutput(item)
		case Fluentd:
			out = makeFluentdOutput(item)
		}

		err := fboClient.Create(context.Background(), &out)
		if err != nil {
			klog.Error(err.Error())
		}
	}
}

func makeElasticsearchOutput(cfg oldConfig) logging.Output {
	var host, prefiex string
	var port int32

	for _, p := range cfg.Params {
		switch p.Name {
		case "Logstash_Prefix":
			prefiex = p.Value
		case "Port":
			portStr := p.Value
			portInt, _ := strconv.ParseInt(portStr, 10, 32)
			port = int32(portInt)
		case "Host":
			host = p.Value
		}
	}

	return logging.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "es",
			Namespace: Namespace,
			Labels: map[string]string{
				"logging.kubesphere.io/enabled":   fmt.Sprint(cfg.Enable),
				"logging.kubesphere.io/component": "logging",
			},
		},
		Spec: logging.OutputSpec{
			Match: "kube.*",
			Elasticsearch: &output.Elasticsearch{
				Host:           host,
				Port:           &port,
				LogstashPrefix: prefiex,
				LogstashFormat: func() *bool {
					b := true
					return &b
				}(),
				TimeKey: "@timestamp",
			},
		},
	}
}

func makeKafkaOutput(cfg oldConfig) logging.Output {
	var brokers, topics string

	for _, p := range cfg.Params {
		switch p.Name {
		case "Brokers":
			brokers = p.Value
		case "Topics":
			topics = p.Value
		}
	}

	return logging.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kafka",
			Namespace: Namespace,
			Labels: map[string]string{
				"logging.kubesphere.io/enabled":   fmt.Sprint(cfg.Enable),
				"logging.kubesphere.io/component": "logging",
			},
		},
		Spec: logging.OutputSpec{
			Match: "kube.*",
			Kafka: &output.Kafka{
				Brokers: brokers,
				Topics:  topics,
			},
		},
	}
}

func makeFluentdOutput(cfg oldConfig) logging.Output {
	var host string
	var port int32

	for _, p := range cfg.Params {
		switch p.Name {
		case "Host":
			host = p.Value
		case "Port":
			portStr := p.Value
			portInt, _ := strconv.ParseInt(portStr, 10, 32)
			port = int32(portInt)
		}
	}

	return logging.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentd",
			Namespace: Namespace,
			Labels: map[string]string{
				"logging.kubesphere.io/enabled":   fmt.Sprint(cfg.Enable),
				"logging.kubesphere.io/component": "logging",
			},
		},
		Spec: logging.OutputSpec{
			Match: "kube.*",
			Forward: &output.Forward{
				Host: host,
				Port: &port,
			},
		},
	}
}
