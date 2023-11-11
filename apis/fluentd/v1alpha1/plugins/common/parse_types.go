package common

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/params"
)

// ParseCommon defines the common parameters for the parse plugin
type ParseCommon struct {
	// The @id parameter specifies a unique name for the configuration.
	Id *string `json:"id,omitempty"`
	// The @type parameter specifies the type of the plugin.
	// +kubebuilder:validation:Enum:=regexp;apache2;apache_error;nginx;syslog;csv;tsv;ltsv;json;multiline;none;grok;multiline_grok
	Type *string `json:"type"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
}

// Parse defines various parameters for the parse plugin
type Parse struct {
	ParseCommon `json:",inline"`
	Time        `json:",inline,omitempty"`
	// Specifies the regular expression for matching logs. Regular expression also supports i and m suffix.
	Expression *string `json:"expression,omitempty"`
	// Specify types for converting field into another, i.e: types user_id:integer,paid:bool,paid_usd_amount:float
	Types *string `json:"types,omitempty"`
	// Specify time field for event time. If the event doesn't have this field, current time is used.
	TimeKey *string `json:"timeKey,omitempty"`
	// If true, use Fluent::Eventnow(current time) as a timestamp when time_key is specified.
	EstimateCurrentEvent *bool `json:"estimateCurrentEvent,omitempty"`
	// If true, keep time field in th record.
	KeepTimeKey *bool `json:"keepTimeKey,omitempty"`
	// Specify timeout for parse processing.
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	Timeout *string `json:"timeout,omitempty"`
	// The pattern of grok.
	GrokPattern *string `json:"grokPattern,omitempty"`
	// Path to the file that includes custom grok patterns.
	CustomPatternPath *string `json:"customPatternPath,omitempty"`
	// The key has grok failure reason.
	GrokFailureKey *string `json:"grokFailureKey,omitempty"`
	// The regexp to match beginning of multiline. This is only for "multiline_grok".
	MultiLineStartRegexp *string `json:"multiLineStartRegexp,omitempty"`
	// Specify grok pattern series set.
	GrokPatternSeries *string `json:"grokPatternSeries,omitempty"`
	// Grok Sections
	Grok []Grok `json:"grok,omitempty"`
}

type Grok struct {
	// The name of this grok section.
	Name *string `json:"name,omitempty"`
	// The pattern of grok. Required parameter.
	Pattern *string `json:"pattern,omitempty"`
	// If true, keep time field in the record.
	KeepTimeKey *bool `json:"keepTimeKey,omitempty"`
	// Specify time field for event time. If the event doesn't have this field, current time is used.
	TimeKey *string `json:"timeKey,omitempty"`
	// Process value using specified format. This is available only when time_type is string
	TimeFormat *string `json:"timeFormat,omitempty"`
	// Use specified timezone. one can parse/format the time value in the specified timezone.
	TimeZone *string `json:"timeZone,omitempty"`
}

func (p *Parse) Name() string {
	return "parse"
}

func (p *Parse) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore("parse")
	if p.Id != nil {
		ps.InsertPairs("@id", fmt.Sprint(*p.Id))
	}
	if p.Type != nil {
		ps.InsertType(fmt.Sprint(*p.Type))
	}
	if p.LogLevel != nil {
		ps.InsertPairs("@log_level", fmt.Sprint(*p.LogLevel))
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
	if p.EstimateCurrentEvent != nil {
		ps.InsertPairs("estimate_current_event", fmt.Sprint(*p.EstimateCurrentEvent))
	}
	if p.KeepTimeKey != nil {
		ps.InsertPairs("keep_time_key", fmt.Sprint(*p.KeepTimeKey))
	}
	if p.Timeout != nil {
		ps.InsertPairs("timeout", fmt.Sprint(*p.Timeout))
	}

	if p.TimeType != nil {
		ps.InsertPairs("time_type", fmt.Sprint(*p.TimeType))
	}
	if p.TimeFormat != nil {
		ps.InsertPairs("time_format", fmt.Sprint(*p.TimeFormat))
	}
	if p.Localtime != nil {
		ps.InsertPairs("localtime", fmt.Sprint(*p.Localtime))
	}
	if p.UTC != nil {
		ps.InsertPairs("utc", fmt.Sprint(*p.UTC))
	}
	if p.Timezone != nil {
		ps.InsertPairs("timezone", fmt.Sprint(*p.Timezone))
	}
	if p.TimeFormatFallbacks != nil {
		ps.InsertPairs("time_format_fallbacks", fmt.Sprint(*p.TimeFormatFallbacks))
	}

	if p.GrokPattern != nil {
		ps.InsertPairs("grok_pattern", fmt.Sprint(*p.GrokPattern))
	}
	if p.CustomPatternPath != nil {
		ps.InsertPairs("custom_pattern_path", fmt.Sprint(*p.CustomPatternPath))
	}
	if p.GrokFailureKey != nil {
		ps.InsertPairs("grok_failure_key", fmt.Sprint(*p.GrokFailureKey))
	}
	if p.MultiLineStartRegexp != nil {
		ps.InsertPairs("multi_line_start_regexp", fmt.Sprint(*p.MultiLineStartRegexp))
	}
	if p.GrokPatternSeries != nil {
		ps.InsertPairs("grok_pattern_series", fmt.Sprint(*p.GrokPatternSeries))
	}
	for _, grok := range p.Grok {
		g := params.NewPluginStore("grok")
		if grok.Name != nil {
			g.InsertPairs("name", fmt.Sprint(*grok.Name))
		}
		if grok.Pattern != nil {
			g.InsertPairs("pattern", fmt.Sprint(*grok.Pattern))
		}
		if grok.KeepTimeKey != nil {
			g.InsertPairs("keep_time_key", fmt.Sprint(*grok.KeepTimeKey))
		}
		if grok.TimeKey != nil {
			g.InsertPairs("time_key", fmt.Sprint(*grok.TimeKey))
		}
		if grok.TimeFormat != nil {
			g.InsertPairs("time_format", fmt.Sprint(*grok.TimeFormat))
		}
		if grok.TimeZone != nil {
			g.InsertPairs("timezone", fmt.Sprint(*grok.TimeZone))
		}
		ps.InsertChilds(g)
	}

	return ps, nil
}
