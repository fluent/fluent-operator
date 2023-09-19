# StatsD

The StatsD input plugin allows you to receive metrics via StatsD protocol.<br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/statsd**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | Listener network interface, default: 0.0.0.0 | string |
| port | UDP port where listening for connections, default: 8125 | *int32 |
