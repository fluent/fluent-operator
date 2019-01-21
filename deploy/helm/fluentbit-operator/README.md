
# FluentBit Operator Chart

[FluentBit Operator](https://github.com/kubesphere/fluentbit-operator) Managed centralized logging component fluent-bit instance on cluster.
## tl;dr:

```bash
$ helm repo add kubesphere http://kubernetes-charts.kubesphere.io/branch/master
$ helm repo update
$ helm install kubesphere/fluentbit-operator
```

## Introduction

This chart bootstraps an [FluentBit Operator](https://github.com/kubesphere/charts/fluentbit-operator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.8+ with Beta APIs enabled

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release kubesphere/fluentbit-operator
```

The command deploys **fluentbit-operator** on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the fluentbit-operator chart and their default values.

|          Parameter          |                Description                               |             Default            |
| --------------------------- | -------------------------------------------------------- | ------------------------------ |
| `image.repository`          | Container image repository                               | `kubesphere/fluentbit-operator` |
| `image.tag       `          | Container image tag                                      | `latest`                       |
| `image.pullPolicy`          | Container pull policy                                    | `Always`                       |
| `tls.enabled`               | Enabled TLS communication between components             | true                           |
| `tls.secretName`            | Specified secret name, which contain tls certs           | This will overwrite automatic Helm certificate generation. |
| `fluentbit.enabled`         | Install fluent-bit                                       | true                           |
| `fluentbit.namespace`       | Specified fluentbit installation namespace               | same as operator namespace     |
| `affinity`                  | Node Affinity                                            | none provided                  |
| `tolerations`               | Node Tolerations                                         | none provided                  |

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```bash
$ helm install --name my-release -f values.yaml kubesphere/fluentbit-operator
```

> **Tip**: You can use the default [values.yaml](values.yaml)


```

