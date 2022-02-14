# Throttle

The Throttle Filter plugin sets the average Rate of messages per Interval, based on leaky bucket and sliding window algorithm. In case of overflood, it will leak within certain rate.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| rate | Amount of messages for the time. | *int64 |
| window | Amount of intervals to calculate average over. | *int64 |
| interval | Time interval, expressed in "sleep" format. e.g 3s, 1.5m, 0.5h etc | string |
| printStatus | Whether to print status messages with current rate and the limits to information logs | *bool |
