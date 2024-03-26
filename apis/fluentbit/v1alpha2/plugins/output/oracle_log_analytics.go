package output

import (
	"fmt"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

type OracleLogAnalytics struct {
	Auth                   *AuthConfig       `json:"auth,omitempty"`
	OCIConfigInRecord      bool              `json:"ociConfigInRecord,omitempty"`
	ObjectStorageNamespace *string           `json:"objectStorageNamespace,omitempty"`
	ProfileName            *string           `json:"profileName,omitempty"`
	LogGroupId             *string           `json:"logGroupId,omitempty"`
	LogSourceName          *string           `json:"logSourceName,omitempty"`
	LogEntityId            *string           `json:"logEntityId,omitempty"`
	LogEntityType          *string           `json:"logEntityType,omitempty"`
	LogPath                *string           `json:"logPath,omitempty"`
	LogSet                 *string           `json:"logSetId,omitempty"`
	GlobalMetadata         map[string]string `json:"globalMetadata,omitempty"`
	LogEventMetadata       map[string]string `json:"logEventMetadata,omitempty"`
	Workers                *int32            `json:"Workers,omitempty"`
	*plugins.TLS           `json:"tls,omitempty"`
}

// +kubebuilder:object:generate:=true

type AuthConfig struct {
	ConfigFileLocation *string `json:"configFileLocation,omitempty"`
	ProfileName        *string `json:"profileName,omitempty"`
}

func (_ *OracleLogAnalytics) Name() string {
	return "oracle_log_analytics"
}

func (o *OracleLogAnalytics) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if o.Auth.ConfigFileLocation != nil {
		kvs.Insert("config_file_location", *o.Auth.ConfigFileLocation)
	}
	if o.Auth.ProfileName != nil {
		kvs.Insert("profile_name", *o.Auth.ProfileName)
	}
	if o.ObjectStorageNamespace != nil {
		kvs.Insert("namespace", *o.ObjectStorageNamespace)
	}
	if o.OCIConfigInRecord {
		kvs.Insert("oci_config_in_record", "true")
	}
	if o.LogGroupId != nil {
		kvs.Insert("oci_la_log_group_id", *o.LogGroupId)
	}
	if o.LogSourceName != nil {
		kvs.Insert("oci_la_log_source_name", *o.LogSourceName)
	}
	if o.LogEntityId != nil {
		kvs.Insert("oci_la_log_entity_id", *o.LogEntityId)
	}
	if o.LogEntityType != nil {
		kvs.Insert("oci_la_log_entity_type", *o.LogEntityType)
	}
	if o.LogPath != nil {
		kvs.Insert("oci_la_log_path", *o.LogPath)
	}
	if o.LogSet != nil {
		kvs.Insert("oci_la_log_set_id", *o.LogSet)
	}
	if o.GlobalMetadata != nil {
		for k, v := range o.GlobalMetadata {
			kvs.Insert("oci_la_global_metadata", fmt.Sprintf("%s %s", k, v))
		}
	}
	if o.LogEventMetadata != nil {
		for k, v := range o.LogEventMetadata {
			kvs.Insert("oci_la_metadata", fmt.Sprintf("%s %s", k, v))
		}
	}
	if o.Workers != nil {
		kvs.Insert("Workers", fmt.Sprint(*o.Workers))
	}
	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
