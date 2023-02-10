# Opensearch

Opensearch defines the parameters for out_opensearch plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | The hostname of your Opensearch node (default: localhost). | *string |
| port | The port number of your Opensearch node (default: 9200). | *uint32 |
| hosts | Hosts defines a list of hosts if you want to connect to more than one Openearch nodes | *string |
| scheme | Specify https if your Opensearch endpoint supports SSL (default: http). | *string |
| path | Path defines the REST API endpoint of Opensearch to post write requests (default: nil). | *string |
| indexName | IndexName defines the placeholder syntax of Fluentd plugin API. See https://docs.fluentd.org/configuration/buffer-section. | *string |
| logstashFormat | If true, Fluentd uses the conventional index name format logstash-%Y.%m.%d (default: false). This option supersedes the index_name option. | *bool |
| logstashPrefix | LogstashPrefix defines the logstash prefix index name to write events when logstash_format is true (default: logstash). | *string |
| user | Optional, The login credentials to connect to Opensearch | *[plugins.Secret](../secret.md) |
| password | Optional, The login credentials to connect to Opensearch | *[plugins.Secret](../secret.md) |
