# Misc
# Path Convention

Path to file in Fluent Bit config should be well regulated. Fluent Bit Operator adopts the following convention internally.

| Dir Path                          | Description                                                  |
| --------------------------------- | ------------------------------------------------------------ |
| /fluent-bit/tail                  | Stores tail related files, eg. file tracking db. Using [fluentbit.spec.positionDB](https://github.com/fluent/fluent-operator/blob/master/docs/fluentbit.md#fluentbitspec) will mount a file `pos.db` under this dir by default. |
| /fluent-bit/secrets/{secret_name} | Stores secrets, eg. TLS files. Specify secrets to mount in [fluentbit.spec.secrets](https://github.com/fluent/fluent-operator/blob/master/docs/fluentbit.md#fluentbitspec), then you have access. |
| /fluent-bit/config                | Stores the main config file and user-defined parser config file. |

> Note that ServiceAccount files are mounted at `/var/run/secrets/kubernetes.io/serviceaccount`.




# Filter Ordering

ClusterFilters run before the namespaced Filters, sorted by `ordinal` and then by name. Set `afterNamespaceFilters: true` on a ClusterFilter to run it after the namespaced Filters instead, e.g. to drop a field that the namespaced Filters still need. `ordinal` orders ClusterFilters within their group, so a ClusterFilter that runs after the namespaced Filters always comes after them, whatever its ordinal. Namespaced Filters reject the field.

```yaml
apiVersion: fluentbit.fluent.io/v1alpha2
kind: ClusterFilter
metadata:
  name: drop-internal-fields
  labels:
    fluentbit.fluent.io/enabled: "true"
spec:
  match: kube.*
  afterNamespaceFilters: true
  filters:
    - recordModifier:
        removeKeys:
          - internal_trace_id
```
