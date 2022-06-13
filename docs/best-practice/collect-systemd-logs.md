# Collect Service Logs of Kubernetes Nodes

This tutorial guides you on how to collect Service logs of kubernetes nodes with fluentbit operator.

# Demo

Firstly, install or update the fluentbit operator

```shell
kubectl apply -f manifests/setup/
kubectl apply -f manifests/logging-stack/fluentbit-fluentBit.yaml
kubectl apply -f manifests/logging-stack/fluentbitconfig-fluentBitConfig.yaml
```

Secondly, change the service logs directory. 
Please create directory `/var/log/journal` if it doesn't exist, and then restart the `systemd-journald` service.

```shell
mkdir /var/log/journal/
systemctl restart systemd-journald
```

Thirdly, set up the fluentbit pipeline. 

```shell
kubectl create cm fluent-bit-lua -n kubesphere-logging-system --from-file=config/scripts/systemd.lua
kubectl apply -f manifests/logging-stack/input-systemd.yaml
kubectl apply -f manifests/logging-stack/filter-systemd.yaml
kubectl apply -f manifests/logging-stack/output-elasticsearch.yaml
```

> This pipeline will send the logs to elasticsearch, it needed a elasticsearch cluster.


> If you want to collect other service logs, such as containerd, you can add a input like the docker input, 
> and modify the systemdFilter.

```bash
    systemdFilter:
      - _SYSTEMD_UNIT=containerd.service
```

For these, the kubelet log will be collected to the elasticsearch. If the fluentbit operator is installed in the 
kubesphere, you can search the log with [Log Search](https://v3-0.docs.fluent.io/docs/toolbox/log-query/).