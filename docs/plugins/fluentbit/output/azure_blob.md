# AzureBlob

Azure Blob is the Azure Blob output plugin, allows to ingest your records into Azure Blob Storage. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/azure_blob**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| accountName | Azure Storage account name | string |
| sharedKey | Specify the Azure Storage Shared Key to authenticate against the storage account | *[plugins.Secret](../secret.md) |
| containerName | Name of the container that will contain the blobs | string |
| blobType | Specify the desired blob type. Must be `appendblob` or `blockblob` | string |
| autoCreateContainer | Creates container if ContainerName is not set. | string |
| path | Optional path to store the blobs. | string |
| emulatorMode | Optional toggle to use an Azure emulator | string |
| endpoint | HTTP Service of the endpoint (if using EmulatorMode) | string |
| compressBlob | Optional: Enables GZIP compression in the final blockblob file. This option isn't compatible when blob_type = appendblob. | string |
| bufferingEnabled | Enable buffering into disk before ingesting into Azure Blob. | *bool |
| uploadFileSize | Specifies the size of files to be uploaded in MB. Defaults to 200M. | string |
| uploadTimeout | Optional. Specify a timeout for uploads. Fluent Bit will start ingesting buffer files which have been created more than x minutes ago and haven't reached upload_file_size limit yet. Defaults to 30m. | string |
| tls | Enable/Disable TLS Encryption. Azure services require TLS to be enabled. | *[plugins.TLS](../tls.md) |
| networking | Include fluentbit networking options for this output-plugin | *plugins.Networking |
