# Elasticsearch

Elasticsearch defines the parameters for out_es output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | The hostname of your Elasticsearch node (default: localhost). | *string |
| port | The port number of your Elasticsearch node (default: 9200). | *uint32 |
| hosts | Hosts defines a list of hosts if you want to connect to more than one Elasticsearch nodes | *string |
| scheme | Specify https if your Elasticsearch endpoint supports SSL (default: http). | *string |
| path | Path defines the REST API endpoint of Elasticsearch to post write requests (default: nil). | *string |
| indexName | IndexName defines the placeholder syntax of Fluentd plugin API. See https://docs.fluentd.org/configuration/buffer-section. | *string |
| logstashFormat | If true, Fluentd uses the conventional index name format logstash-%Y.%m.%d (default: false). This option supersedes the index_name option. | *bool |
| logstashPrefix | LogstashPrefix defines the logstash prefix index name to write events when logstash_format is true (default: logstash). | *string |
| user | Optional, The login credentials to connect to Elasticsearch | *[plugins.Secret](../secret.md) |
| password | Optional, The login credentials to connect to Elasticsearch | *[plugins.Secret](../secret.md) |
| sslVerify | Optional, Force certificate validation | *bool |
| caFile | Optional, Absolute path to CA certificate file | *string |
| clientCert | Optional, Absolute path to client Certificate file | *string |
| clientKey | Optional, Absolute path to client private Key file | *string |
| clientKeyPassword | Optional, password for ClientKey file | *[plugins.Secret](../secret.md) |
