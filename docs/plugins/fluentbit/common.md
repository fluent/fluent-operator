# CommonParams




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| alias | Alias for the plugin | string |
| retryLimit | RetryLimit describes how many times fluent-bit should retry to send data to a specific output. If set to false fluent-bit will try indefinetly. If set to any integer N>0 it will try at most N+1 times. Leading zeros are not allowed (values such as 007, 0150, 01 do not work). If this property is not defined fluent-bit will use the default value: 1. | string |
