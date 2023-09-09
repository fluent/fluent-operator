# MQTT

The MQTT input plugin, allows to retrieve messages/data from MQTT control packets over a TCP connection. <br /> The incoming data to receive must be a JSON map. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/mqtt**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | Listener network interface, default: 0.0.0.0 | string |
| port | TCP port where listening for connections, default: 1883 | *int32 |
