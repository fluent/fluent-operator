# AzureLogsIngestion

Azure Logs Ingestion is the Azure Logs Ingestion output plugin, allows you to ingest your records into Azure Monitor tables through a data collection endpoint and rule. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/azure_logs_ingestion**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tenantID | Azure Active Directory tenant ID | string |
| clientID | Azure Active Directory application (client) ID | string |
| clientSecret | Azure Active Directory application client secret | *[plugins.Secret](../secret.md) |
| dceURL | Data Collection Endpoint (DCE) URL | string |
| dcrID | Data Collection Rule (DCR) immutable ID | string |
| tableName | Name of the destination table the records are ingested into | string |
| timeKey | Specify the name of the key where the timestamp is stored. | string |
| timeGenerated | If set, overrides the timeKey value with the `time-generated-field` HTTP header value. | *bool |
| compress | Enable gzip compression of the request body. | *bool |
