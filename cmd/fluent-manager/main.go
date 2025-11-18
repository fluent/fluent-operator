/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"

	"github.com/fluent/fluent-operator/v3/controllers"
	// +kubebuilder:scaffold:imports
)

const (
	fluentBitName = "fluent-bit"
	fluentdName   = "fluentd"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var logPath string
	var tlsOpts []func(*tls.Config)

	zapOpts := &zap.Options{Development: true}
	opts := NewOptions(zapOpts)
	opts.ParseFlags()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(zapOpts)))

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	if !opts.EnableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	// Create watchers for metrics and webhooks certificates
	var metricsCertWatcher, webhookCertWatcher *certwatcher.CertWatcher

	// Initial webhook TLS options
	webhookTLSOpts := tlsOpts
	if len(opts.WebhookCertPath) > 0 {
		setupLog.Info(
			"Initializing webhook certificate watcher using provided certificates",
			"webhook-cert-path", opts.WebhookCertPath,
			"webhook-cert-name", opts.WebhookCertName,
			"webhook-cert-key", opts.WebhookCertKey,
		)

		var err error
		webhookCertWatcher, err = certwatcher.New(
			filepath.Join(opts.WebhookCertPath, opts.WebhookCertName),
			filepath.Join(opts.WebhookCertPath, opts.WebhookCertKey),
		)
		if err != nil {
			setupLog.Error(err, "Failed to initialize webhook certificate watcher")
			os.Exit(1)
		}

		webhookTLSOpts = append(webhookTLSOpts, func(config *tls.Config) {
			config.GetCertificate = webhookCertWatcher.GetCertificate
		})
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: webhookTLSOpts,
	})

	// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
	// More info:
	// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/metrics/server
	// - https://book.kubebuilder.io/reference/metrics.html
	metricsServerOptions := metricsserver.Options{
		BindAddress:   opts.MetricsAddr,
		SecureServing: opts.SecureMetrics,
		TLSOpts:       tlsOpts,
	}

	if opts.SecureMetrics {
		// FilterProvider is used to protect the metrics endpoint with authn/authz.
		// These configurations ensure that only authorized users and service accounts
		// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
		// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/metrics/filters#WithAuthenticationAndAuthorization
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}
	// If the certificate is not specified, controller-runtime will automatically
	// generate self-signed certificates for the metrics server. While convenient for development and testing,
	// this setup is not recommended for production.
	//
	// TODO(user): If you enable certManager, uncomment the following lines:
	// - [METRICS-WITH-CERTS] at config/default/kustomization.yaml to generate and use certificates
	// managed by cert-manager for the metrics server.
	// - [PROMETHEUS-WITH-CERTS] at config/prometheus/kustomization.yaml for TLS certification.
	if len(opts.MetricsCertPath) > 0 {
		setupLog.Info(
			"Initializing metrics certificate watcher using provided certificates",
			"metrics-cert-path", opts.MetricsCertPath,
			"metrics-cert-name", opts.MetricsCertName,
			"metrics-cert-key", opts.MetricsCertKey,
		)

		var err error
		metricsCertWatcher, err = certwatcher.New(
			filepath.Join(opts.MetricsCertPath, opts.MetricsCertName),
			filepath.Join(opts.MetricsCertPath, opts.MetricsCertKey),
		)
		if err != nil {
			setupLog.Error(err, "to initialize metrics certificate watcher", "error", err)
			os.Exit(1)
		}

		metricsServerOptions.TLSOpts = append(metricsServerOptions.TLSOpts, func(config *tls.Config) {
			config.GetCertificate = metricsCertWatcher.GetCertificate
		})
	}

	ctrlOpts := ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsServerOptions,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: opts.ProbeAddr,
		LeaderElection:         opts.EnableLeaderElection,
		LeaderElectionID:       "45c4fdd2.fluent.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	}

	namespacedController := false
	if opts.WatchNamespaces != "" {
		config := cache.Config{}
		namespacedController = true

		ctrlOpts.Cache.DefaultNamespaces = make(map[string]cache.Config)
		for namespace := range strings.SplitSeq(opts.WatchNamespaces, ",") {
			ctrlOpts.Cache.DefaultNamespaces[namespace] = config
		}
	}

	logPath = getContainerLogPath("/fluent-operator/fluent-bit.env")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrlOpts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	fluentBitEnabled, fluentdEnabled := true, true
	if opts.DisabledControllers != "" {
		switch opts.DisabledControllers {
		case fluentBitName:
			fluentBitEnabled = false
		case fluentdName:
			fluentdEnabled = false
		default:
			err := errors.New("incorrect value for disable-component-controllers flag")
			setupLog.Error(err, "it will not be proceeded (possible values are: fluent-bit, fluentd)")
		}
	}

	if fluentBitEnabled {
		utilruntime.Must(fluentbitv1alpha2.AddToScheme(scheme))
		if err = (&controllers.FluentBitConfigReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("FluentBitConfig"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "FluentBitConfig")
			os.Exit(1)
		}

		if err = (&controllers.CollectorReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Collector"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Collector")
			os.Exit(1)
		}

		if err = (&controllers.FluentBitReconciler{
			Client:               mgr.GetClient(),
			Log:                  ctrl.Log.WithName("controllers").WithName("FluentBit"),
			Scheme:               mgr.GetScheme(),
			ContainerLogRealPath: logPath,
			Namespaced:           namespacedController,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "FluentBit")
			os.Exit(1)
		}

		// +kubebuilder:scaffold:builder
	}

	if fluentdEnabled {
		utilruntime.Must(fluentdv1alpha1.AddToScheme(scheme))
		if err = (&controllers.FluentdConfigReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("FluentdConfig"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "FluentdConfig")
			os.Exit(1)
		}

		if err = (&controllers.FluentdReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Fluentd"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Fluentd")
			os.Exit(1)
		}

		// +kubebuilder:scaffold:builder
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// getContainerLogPath determines the container log path for FluentBit.
func getContainerLogPath(envFilePath string) string {
	var logPath string

	if envLogPath := os.Getenv("CONTAINER_LOG_PATH"); envLogPath != "" {
		logPath = envLogPath
	} else if envs, err := godotenv.Read(envFilePath); err == nil {
		logPath = envs["CONTAINER_ROOT_DIR"] + "/containers"
	}

	if logPath == "" {
		logPath = "/var/log/containers"
	}

	return logPath
}
