package common

import (
	"fmt"

	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/params"
)

type ParseCommon struct {
	// The @id parameter specifies a unique name for the configuration.
	Id *string `json:"id,omitempty"`
	// The @type parameter specifies the type of the plugin.
	// +kubebuilder:validation:Enum:=regexp;apache2;apache_error;nginx;syslog;csv;tsv;ltsv;json;multiline;none
	Type *string `json:"type"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
}

type Parse struct {
	*ParseCommon `json:",inline"`
	*Time        `json:",inline,omitempty"`
	// Specifies the regular expression for matching logs. Regular expression also supports i and m suffix.
	Expression *string `json:"expression,omitempty"`
	// Specify types for converting field into another, i.e: types user_id:integer,paid:bool,paid_usd_amount:float
	Types *string `json:"types,omitempty"`
	// Specify time field for event time. If the event doesn't have this field, current time is used.
	TimeKey *string `json:"timeKey,omitempty"`
	// If true, use Fluent::EventTime.now(current time) as a timestamp when time_key is specified.
	EstimateCurentEvent *bool `json:"estimateCurrentEvent,omitempty"`
	// If true, keep time field in th record.
	KeepTimeKey *bool `json:"keepTimeKey,omitempty"`
	// Specify timeout for parse processing.
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	Timeout *string `json:"timeout,omitempty"`
}

func (p *Parse) Name() string {
	return "parse"
}

func (p *Parse) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore("parse")
	if p.ParseCommon.Id != nil {
		ps.InsertPairs("@id", fmt.Sprint(*p.ParseCommon.Id))
	}
	if p.ParseCommon.Type != nil {
		ps.InsertPairs("@type", fmt.Sprint(*p.ParseCommon.Type))
	}
	if p.ParseCommon.LogLevel != nil {
		ps.InsertPairs("@log_level", fmt.Sprint(*p.ParseCommon.LogLevel))
	}
	if p.Expression != nil {
		ps.InsertPairs("expression", fmt.Sprint(*p.Expression))
	}
	if p.Types != nil {
		ps.InsertPairs("types", fmt.Sprint(*p.Types))
	}
	if p.TimeKey != nil {
		ps.InsertPairs("time_key", fmt.Sprint(*p.TimeKey))
	}
	if p.EstimateCurentEvent != nil {
		ps.InsertPairs("estimate_curent_event", fmt.Sprint(*p.EstimateCurentEvent))
	}
	if p.KeepTimeKey != nil {
		ps.InsertPairs("keep_timeout", fmt.Sprint(*p.KeepTimeKey))
	}
	if p.Timeout != nil {
		ps.InsertPairs("timeout", fmt.Sprint(*p.Timeout))
	}
	if p.Time != nil {
		if p.Time.TimeType != nil {
			ps.InsertPairs("time_type", fmt.Sprint(*p.Time.TimeType))
		}
		if p.Time.TimeFormat != nil {
			ps.InsertPairs("time_type", fmt.Sprint(*p.Time.TimeFormat))
		}
		if p.Time.Localtime != nil {
			ps.InsertPairs("localtime", fmt.Sprint(*p.Time.Localtime))
		}
		if p.Time.UTC != nil {
			ps.InsertPairs("utc", fmt.Sprint(*p.Time.UTC))
		}
		if p.Time.Timezone != nil {
			ps.InsertPairs("timezone", fmt.Sprint(*p.Time.Timezone))
		}
		if p.Time.TimeFormatFallbacks != nil {
			ps.InsertPairs("time_format_fallbacks", fmt.Sprint(*p.Time.TimeFormatFallbacks))
		}
	}
	return ps, nil
}
