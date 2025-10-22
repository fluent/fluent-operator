package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Azure Blob is the Azure Blob output plugin, allows to ingest your records into Azure Blob Storage. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/azure_blob**
type AzureBlob struct {
	// Azure Storage account name
	AccountName string `json:"accountName"`
	// Specify the Azure Storage Shared Key to authenticate against the storage account
	SharedKey *plugins.Secret `json:"sharedKey"`
	// Name of the container that will contain the blobs
	ContainerName string `json:"containerName"`
	// Specify the desired blob type. Must be `appendblob` or `blockblob`
	// +kubebuilder:validation:Enum:=appendblob;blockblob
	BlobType string `json:"blobType,omitempty"`
	// Creates container if ContainerName is not set.
	// +kubebuilder:validation:Enum:=on;off
	AutoCreateContainer string `json:"autoCreateContainer,omitempty"`
	// Optional path to store the blobs.
	Path string `json:"path,omitempty"`
	// Optional toggle to use an Azure emulator
	// +kubebuilder:validation:Enum:=on;off
	EmulatorMode string `json:"emulatorMode,omitempty"`
	// HTTP Service of the endpoint (if using EmulatorMode)
	Endpoint string `json:"endpoint,omitempty"`
	// Optional: Enables GZIP compression in the final blockblob file. This option isn't compatible when blob_type = appendblob.
	// +kubebuilder:validation:Enum:=on;off
	CompressBlob string `json:"compressBlob,omitempty"`
	// Enable buffering into disk before ingesting into Azure Blob.
	BufferingEnabled *bool `json:"bufferingEnabled,omitempty"`
	// Specifies the size of files to be uploaded in MB. Defaults to 200M.
	// +kubebuilder:default:="200M"
	UploadFileSize string `json:"uploadFileSize,omitempty"`
	// Optional. Specify a timeout for uploads. Fluent Bit will start ingesting buffer files which have been created more than x minutes ago and haven't reached upload_file_size limit yet. Defaults to 30m.
	// +kubebuilder:default:="30m"
	UploadTimeout string `json:"uploadTimeout,omitempty"`
	// Enable/Disable TLS Encryption. Azure services require TLS to be enabled.
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

// Name implement Section() method
func (*AzureBlob) Name() string {
	return "azure_blob"
}

// Params implement Section() method
func (o *AzureBlob) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "shared_key", o.SharedKey, sl); err != nil {
		return nil, err
	}

	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if o.Networking != nil {
		net, err := o.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "account_name", o.AccountName)
	plugins.InsertKVString(kvs, "container_name", o.ContainerName)
	plugins.InsertKVString(kvs, "blob_type", o.BlobType)
	plugins.InsertKVString(kvs, "auto_create_container", o.AutoCreateContainer)
	plugins.InsertKVString(kvs, "compress_blob", o.CompressBlob)
	plugins.InsertKVField(kvs, "buffering_enabled", o.BufferingEnabled)
	plugins.InsertKVString(kvs, "upload_file_size", o.UploadFileSize)
	plugins.InsertKVString(kvs, "upload_timeout", o.UploadTimeout)
	plugins.InsertKVString(kvs, "path", o.Path)
	plugins.InsertKVString(kvs, "emulator_mode", o.EmulatorMode)
	plugins.InsertKVString(kvs, "endpoint", o.Endpoint)

	return kvs, nil
}
