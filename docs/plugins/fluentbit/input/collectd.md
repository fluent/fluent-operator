# Collectd

The Collectd input plugin allows you to receive datagrams from collectd service. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/collectd**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | Set the address to listen to, default: 0.0.0.0 | string |
| port | Set the port to listen to, default: 25826 | *int32 |
| typesDB | Set the data specification file,default: /usr/share/collectd/types.db | string |
