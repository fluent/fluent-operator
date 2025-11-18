package main

import (
	"flag"

	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type Options struct {
	WatchNamespaces      string
	MetricsAddr          string
	EnableLeaderElection bool
	SecureMetrics        bool
	WebhookCertPath      string
	WebhookCertName      string
	WebhookCertKey       string
	MetricsCertPath      string
	MetricsCertName      string
	MetricsCertKey       string
	EnableHTTP2          bool
	ProbeAddr            string
	DisabledControllers  string
}

func NewOptions(zapOpts *zap.Options) *Options {
	opts := new(Options)
	flag.StringVar(&opts.WatchNamespaces, "watch-namespaces", "",
		"Optional comma separated list of namespaces to watch for resources in. Defaults to cluster scope.")
	flag.StringVar(&opts.MetricsAddr, "metrics-bind-address", "0",
		"The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP, or leave "+
			"as 0 to disable the metrics service.")
	flag.BoolVar(&opts.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&opts.SecureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.StringVar(&opts.WebhookCertPath, "webhook-cert-path", "",
		"The directory that contains the webhook certificate.")
	flag.StringVar(&opts.WebhookCertName, "webhook-cert-name", "tls.crt",
		"The name of the webhook certificate file.")
	flag.StringVar(&opts.WebhookCertKey, "webhook-cert-key", "tls.key", "The name of the webhook key file.")
	flag.StringVar(&opts.MetricsCertPath, "metrics-cert-path", "",
		"The directory that contains the metrics server certificate.")
	flag.StringVar(&opts.MetricsCertName, "metrics-cert-name", "tls.crt",
		"The name of the metrics server certificate file.")
	flag.StringVar(&opts.MetricsCertKey, "metrics-cert-key", "tls.key", "The name of the metrics server key file.")
	flag.BoolVar(&opts.EnableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	flag.StringVar(&opts.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&opts.DisabledControllers, "disable-component-controllers", "",
		"Optional argument that accepts two values: fluent-bit and fluentd. "+
			"The specific controller will not be started if it's disabled.")

	zapOpts.BindFlags(flag.CommandLine)

	return opts
}

func (o *Options) ParseFlags() {
	flag.Parse()
}
