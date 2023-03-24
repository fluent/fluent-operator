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
	"flag"
	"os"
	"strings"

	"errors"

	"github.com/joho/godotenv"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"

	"github.com/fluent/fluent-operator/v2/controllers"
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
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var watchNamespaces string
	var logPath string
	var disabledControllers string
	flag.StringVar(&watchNamespaces, "watch-namespaces", "", "Optional comma separated list of namespaces to watch for resources in. Defaults to cluster scope.")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&disabledControllers, "disable-component-controllers", "",
		"Optional argument that accepts two values: fluent-bit and fluentd. "+
			"The specific controller will not be started if it's disabled.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	if envs, err := godotenv.Read("/fluent-operator/fluent-bit.env"); err == nil {
		logPath = envs["CONTAINER_ROOT_DIR"] + "/containers"
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	ctrlOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "45c4fdd2.fluent.io",
	}
	namespacedController := false
	if watchNamespaces != "" {
		namespacedController = true
		namespaces := strings.Split(watchNamespaces, ",")
		if len(namespaces) > 1 {
			ctrlOpts.NewCache = cache.MultiNamespacedCacheBuilder(namespaces)
		} else {
			ctrlOpts.Namespace = namespaces[0]
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrlOpts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	fluentBitEnabled, fluentdEnabled := true, true
	if disabledControllers != "" {
		if disabledControllers == fluentBitName {
			fluentBitEnabled = false
		} else if disabledControllers == fluentdName {
			fluentdEnabled = false
		} else {
			setupLog.Error(errors.New("incorrect value for `-disable-component-controllers` and it will not be proceeded (possible values are: fluent-bit, fluentd)"), "")
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
