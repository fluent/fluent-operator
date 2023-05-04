# Secret

Secret defines the key of a value.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| valueFrom |  | ValueSource |
# ValueSource

ValueSource defines how to find a value's key.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| secretKeyRef | Selects a key of a secret in the pod's namespace | [corev1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core) |
# SecretLoader




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| Client |  | client.Client |
