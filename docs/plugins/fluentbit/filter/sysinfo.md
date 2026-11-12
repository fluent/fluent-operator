# Sysinfo

The Sysinfo Filter plugin appends system information, such as the Fluent Bit version or the hostname, to each record. <br /> Only the keys that are set are appended. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/sysinfo**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| fluentBitVersionKey | Key to append the Fluent Bit version under. | string |
| osNameKey | Key to append the operating system name under, e.g. linux. | string |
| hostnameKey | Key to append the hostname under. In Kubernetes this is the hostname of the Fluent Bit pod, which is the pod name unless hostNetwork is enabled, in which case it is the node hostname. | string |
| osVersionKey | Key to append the operating system version under. Only supported on Linux. | string |
| kernelVersionKey | Key to append the kernel version under. Only supported on Linux. | string |
