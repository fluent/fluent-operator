# Sample

The in_sample input plugin generates sample events. It is useful for testing, debugging, benchmarking and getting started with Fluentd.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | The tag of the event. The value is the tag assigned to the generated events. | *string |
| size | The number of events in the event stream of each emit. | *int64 |
| rate | It configures how many events to generate per second. | *int64 |
| autoIncrementKey | If specified, each generated event has an auto-incremented key field. | *string |
| sample | The sample data to be generated. It should be either an array of JSON hashes or a single JSON hash. If it is an array of JSON hashes, the hashes in the array are cycled through in order. | *string |
