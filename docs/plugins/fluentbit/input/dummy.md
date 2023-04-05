# Dummy

The dummy input plugin, generates dummy events. <br /> It is useful for testing, debugging, benchmarking and getting started with Fluent Bit. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/dummy**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all records comming from this plugin. | string |
| dummy | Dummy JSON record. | string |
| rate | Events number generated per second. | *int32 |
| samples | Sample events to generate. | *int32 |
