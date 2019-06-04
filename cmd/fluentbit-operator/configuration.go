package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"kubesphere.io/fluentbit-operator/cmd/fluentbit-operator/fluentbit"
	"path/filepath"
)

//Initialize the configuration
const configFile = "/fluentbit-operator/config/config.toml"

// Init the configuration
func Init() {
	logrus.Info("Initializing configuration")
	viper.SetDefault("tls.enabled", false)
	viper.SetDefault("tls.sharedKey", "Thei6pahshubajee")
	viper.SetDefault("fluent-bit.image", "dockerhub.qingcloud.com/kslogging/fluent-bit:1.0.4")
	viper.SetDefault("fluent-bit.containersLogMountedPath", "/var/lib/docker/containers")
	viper.SetDefault("fluent-bit.pullPolicy", corev1.PullIfNotPresent)
	viper.SetDefault("configmap-reload.image", "dockerhub.qingcloud.com/kslogging/configmap-reload:latest")
	go handleConfigChanges()
}

func handleConfigChanges() {
	c := make(chan fsnotify.Event, 1)
	viper.SetConfigFile(configFile)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			logrus.Fatal(err)
		}
		defer watcher.Close()

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		configFile := filepath.Clean(configFile)
		configDir, _ := filepath.Split(configFile)

		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					// we only care about the config file or the ConfigMap directory (if in Kubernetes)
					if filepath.Clean(event.Name) == configFile || filepath.Base(event.Name) == "..data" {
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
							err := viper.ReadInConfig()
							if err != nil {
								logrus.Println("error:", err)
							}
							c <- event
						}
					}
				case err := <-watcher.Errors:
					logrus.Println("error:", err)
				}
			}
		}()

		watcher.Add(configDir)
		<-done
	}()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error during reading config file : %s", err))
	}
	c <- fsnotify.Event{Name: "Initial", Op: fsnotify.Create}

	for e := range c {
		logrus.Infoln("New config file change", e.String())
		configureOperator()
	}
}

func configureOperator() {
	if viper.GetBool("fluent-bit.enabled") {
		logrus.Info("Trying to init fluent-bit")
		fluentbit.InitFluentBit(GlobalLabels)
	} else if !viper.GetBool("fluent-bit.enabled") {
		logrus.Info("Deleting fluent-bit DaemonSet...")
		fluentbit.DeleteFluentBit(GlobalLabels)
	}
}
