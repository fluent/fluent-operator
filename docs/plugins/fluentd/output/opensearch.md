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
| sslVerify | Optional, Force certificate validation | *bool |
| caFile | Optional, Absolute path to CA certificate file | *string |
| clientCert | Optional, Absolute path to client Certificate file | *string |
| clientKey | Optional, Absolute path to client private Key file | *string |
| clientKeyPassword | Optional, password for ClientKey file | *[plugins.Secret](../secret.md) |
| sslVersion | Optional, You can specify SSL/TLS version (default: TLSv1_2) | *string |
| sslMinVersion | Optional, Minimum SSL/TLS version | *string |
| sslMaxVersion | Optional, Maximum SSL/TLS version | *string |
| logOs400Reason | Optional, Enable logging of 400 reason without enabling debug log level (default: false) | *bool |
| requestTimeout | Optional, HTTP request timeout in seconds (default: 5s) | *string |
| reconnectOnError | Optional, Indicates that the plugin should reset connection on any error (reconnect on next send) (default: false) | *bool |
| reloadConnections | Optional, Automatically reload connection after 10000 documents (default: true) | *bool |
| reloadAfter | Optional, When ReloadConnections true, this is the integer number of operations after which the plugin will reload the connections (default: 10000) | *uint32 |
| reloadOnFailure | Optional, Indicates that the opensearch-transport will try to reload the nodes addresses if there is a failure while making the request (default: false) | *bool |
| maxRetryGetOsVersion | Optional, You can specify times of retry obtaining OpenSearch version (default: 15) | *uint32 |
| failOnDetectingOsVersionRetryExceed | Optional, Indicates whether to fail when max_retry_get_os_version is exceeded (default: true) | *bool |
| defaultOpensearchVersion | Optional, Default OpenSearch version (default: 1) | *uint32 |
| verifyOsVersionAtStartup | Optional, Validate OpenSearch version at startup (default: true) | *bool |
| templateOverwrite | Optional, Always update the template, even if it already exists (default: false) | *bool |
| maxRetryPuttingTemplate | Optional, You can specify times of retry putting template (default: 10) | *uint32 |
| failOnPuttingTemplateRetryExceed | Optional, Indicates whether to fail when max_retry_putting_template is exceeded (default: true) | *bool |
| snifferClassName | Optional, Provide a different sniffer class name | *string |
| selectorClassName | Optional, Provide a selector class name | *string |
| httpBackend | Optional, You can specify HTTP backend (default: excon). Options: excon, typhoeus | *string |
| httpBackendExconNonblock | Optional, With http_backend_excon_nonblock false, plugin uses excon with nonblock=false (default: true) | *bool |
| compressionLevel | Optional, You can specify the compression level (default: no_compression). Options: no_compression, best_compression, best_speed, default_compression | *string |
| preferOjSerializer | Optional, With default behavior, plugin uses Yajl as JSON encoder/decoder. Set to true to use Oj (default: false) | *bool |
| suppressTypeName | Optional, Suppress '[types removal]' warnings on OpenSearch 2.x (default: true for OS2+) | *bool |
| contentType | Optional, With content_type application/x-ndjson, plugin adds application/x-ndjson as Content-Type (default: application/json) | *string |
| includeTagKey | Optional, Include tag key in record (default: false) | *bool |
| tagKey | Optional, Tag key name when include_tag_key is true (default: tag) | *string |
| idKey | Optional, Record accessor syntax to specify the field to use as _id in OpenSearch | *string |
| removeKeys | Optional, Remove specified keys from the event record | *string |
| removeKeysOnUpdate | Optional, Remove keys when record is being updated | *string |
| writeOperation | Optional, The write operation (default: index). Options: index, create, update, upsert | *string |
| emitErrorForMissingId | Optional, When write_operation is not index, setting this true will cause plugin to emit_error_event of records which do not include _id field (default: false) | *bool |
| customHeaders | Optional, Custom headers in Hash format | *string |
| pipeline | Optional, Pipeline name | *string |
| utcIndex | Optional, UTC index (default: false for local time) | *bool |
| suppressDocWrap | Optional, Suppress doc_wrap (default: false) | *bool |
| ignoreExceptions | Optional, List of exception classes to ignore | *string |
| exceptionBackup | Optional, Backup chunk when ignore exception occurs (default: true) | *bool |
| bulkMessageRequestThreshold | Optional, Configure bulk_message request splitting threshold size (default: -1 unlimited) | *int32 |
| applicationName | Optional, Specify the application name for the rollover index to be created (default: default) | *string |
| indexDatePattern | Optional, Specify the index date pattern for creating a rollover index (default: now/d) | *string |
| useLegacyTemplate | Optional, Use legacy template or not (default: false for composable templates) | *bool |
