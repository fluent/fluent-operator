# Collect Service Logs of Kubernetes Nodes

This tutorial guides you on how to collect Service logs of kubernetes nodes with fluentbit operator.

# Demo

First, install or update the fluentbit operator

```shell
kubectl apply -f manifests/setup/
kubectl apply -f manifests/logging-stack/fluentbit-fluentBit.yaml
kubectl apply -f manifests/logging-stack/fluentbitconfig-fluentBitConfig.yaml
```

Second, change the service logs directory

```shell
mkdir /var/log/journal/
systemctl restart systemd-journald
```

Three, set up the fluentbit pipeline. 

```shell
kubectl create cm fluent-bit-lua -n kubesphere-logging-system --from-file=config/scripts/systemd.lua
kubectl apply -f manifests/logging-stack/input-systemd.yaml
kubectl apply -f manifests/logging-stack/filter-systemd.yaml
kubectl apply -f manifests/logging-stack/output-systemd-elasticsearchyaml
```

Note: This pipeline will send the logs to elasticsearch, it needed a elasticsearch cluster.

For these, the kubelet and docker log will be collected to the elasticsearch. If the fluentbit operator is installed in the 
kubesphere, you can search the log with [Log Search](https://v3-0.docs.kubesphere.io/docs/toolbox/log-query/)