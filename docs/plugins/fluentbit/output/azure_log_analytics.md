# AzureLogAnalytics

Azure Log Analytics is the Azure Log Analytics output plugin, allows you to ingest your records into Azure Log Analytics Workspace. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/azure**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| customerID | Customer ID or Workspace ID | *[plugins.Secret](../secret.md) |
| sharedKey | Specify the primary or the secondary client authentication key | *[plugins.Secret](../secret.md) |
| logType | Name of the event type. | string |
| timeKey | Specify the name of the key where the timestamp is stored. | string |
| timeGenerated | If set, overrides the timeKey value with the `time-generated-field` HTTP header value. | *bool |
