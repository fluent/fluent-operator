# Secret




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| valueFrom | ValueSource represents a source for the value of a secret. | ValueSource |
# ValueSource




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| secretKeyRef | Selects a key of a secret in the pod's namespace | [corev1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core) |
