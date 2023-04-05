# Throttle

Throttle filter allows you to set the average rate of messages per internal, based on leaky bucket and sliding window algorithm. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/throttle**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| rate | Rate is the amount of messages for the time. | *int64 |
| window | Window is the amount of intervals to calculate average over. | *int64 |
| interval | Interval is the time interval expressed in \"sleep\" format. e.g. 3s, 1.5m, 0.5h, etc. | string |
| printStatus | PrintStatus represents whether to print status messages with current rate and the limits to information logs. | *bool |
