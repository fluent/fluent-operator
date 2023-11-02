package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/custom"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/params"
)

// InputCommon defines the common parameters for input plugins
type InputCommon struct {
	// The @id parameter specifies a unique name for the configuration.
	Id *string `json:"id,omitempty"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
	// The @label parameter is to route the input events to <label> sections.
	Label *string `json:"label,omitempty"`
}

// Input defines all available input plugins and their parameters
type Input struct {
	InputCommon `json:",inline"`
	// in_forward plugin
	Forward *Forward `json:"forward,omitempty"`
	// in_http plugin
	Http *Http `json:"http,omitempty"`
	// in_tail plugin
	Tail *Tail `json:"tail,omitempty"`
	// in_sample plugin
	Sample *Sample `json:"sample,omitempty"`
	// Custom plugin type
	CustomPlugin *custom.CustomPlugin `json:"customPlugin,omitempty"`
	// monitor_agent plugin
	MonitorAgent *MonitorAgent `json:"monitorAgent,omitempty"`
}

// DeepCopyInto implements the DeepCopyInto interface.
func (in *Input) DeepCopyInto(out *Input) {
	bytes, err := json.Marshal(*in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		panic(err)
	}
}

func (i *Input) Name() string {
	return "source"
}

func (i *Input) Params(loader plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(i.Name())

	if i.Id != nil {
		ps.InsertPairs("@id", fmt.Sprint(*i.Id))
	}

	if i.LogLevel != nil {
		ps.InsertPairs("@log_level", fmt.Sprint(*i.LogLevel))
	}

	if i.Label != nil {
		ps.InsertPairs("@label", fmt.Sprint(*i.Label))
	}

	if i.Forward != nil {
		ps.InsertType(string(params.ForwardInputType))
		return i.forwardPlugin(ps, loader), nil
	}

	if i.Http != nil {
		ps.InsertType(string(params.HttpInputType))
		return i.httpPlugin(ps, loader), nil
	}

	if i.Tail != nil {
		ps.InsertType(string(params.TailInputType))
		return i.tailPlugin(ps, loader), nil
	}

	if i.Sample != nil {
		ps.InsertType(string(params.SampleInputType))
		return i.samplePlugin(ps, loader), nil
	}

	if i.CustomPlugin != nil {
		customPs, _ := i.CustomPlugin.Params(loader)
		ps.Content = customPs.Content
		return ps, nil
	}

	if i.MonitorAgent != nil {
		ps.InsertType(string(params.MonitorAgentType))
		return i.monitorAgentPlugin(ps, loader), nil
	}

	return nil, errors.New("you must define an input plugin")
}

func (i *Input) tailPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	tailModel := i.Tail
	childs := make([]*params.PluginStore, 0)

	if tailModel.Parse != nil {
		child, _ := tailModel.Parse.Params(loader)
		childs = append(childs, child)
	}

	if tailModel.Group != nil {
		child, _ := tailModel.Group.Params(loader)
		childs = append(childs, child)
	}
	// TODO: add group section!
	parent.InsertChilds(childs...)

	if tailModel.Tag != "" {
		parent.InsertPairs("tag", fmt.Sprint(tailModel.Tag))
	}

	if tailModel.Path != "" {
		parent.InsertPairs("path", fmt.Sprint(tailModel.Path))
	}

	if tailModel.PathTimezone != "" {
		parent.InsertPairs("path_timezone", fmt.Sprint(tailModel.PathTimezone))
	}

	if tailModel.ExcludePath != nil {
		parent.InsertPairs("exclude_path", strings.ReplaceAll(fmt.Sprintf("%+q", tailModel.ExcludePath), "\" \"", "\", \""))
	}

	if tailModel.FollowInodes != nil {
		parent.InsertPairs("follow_inodes", fmt.Sprint(*tailModel.FollowInodes))
	}

	if tailModel.RefreshInterval != nil {
		parent.InsertPairs("refresh_interval", fmt.Sprint(*tailModel.RefreshInterval))
	}

	if tailModel.LimitRecentlyModified != nil {
		parent.InsertPairs("limit_recently_modified", fmt.Sprint(*tailModel.LimitRecentlyModified))
	}

	if tailModel.SkipRefreshOnStartup != nil {
		parent.InsertPairs("skip_refresh_on_startup", fmt.Sprint(*tailModel.SkipRefreshOnStartup))
	}

	if tailModel.ReadFromHead != nil {
		parent.InsertPairs("read_from_head", fmt.Sprint(*tailModel.ReadFromHead))
	}

	if tailModel.Encoding != "" {
		parent.InsertPairs("encoding", fmt.Sprint(tailModel.Encoding))
	}

	if tailModel.FromEncoding != "" {
		parent.InsertPairs("from_encoding", fmt.Sprint(tailModel.FromEncoding))
	}

	if tailModel.ReadLinesLimit != nil {
		parent.InsertPairs("read_lines_limit", fmt.Sprint(*tailModel.ReadLinesLimit))
	}

	if tailModel.ReadBytesLimitPerSecond != nil {
		parent.InsertPairs("read_bytes_limit_per_second", fmt.Sprint(*tailModel.ReadBytesLimitPerSecond))
	}

	if tailModel.MaxLineSize != nil {
		parent.InsertPairs("max_line_size", fmt.Sprint(*tailModel.MaxLineSize))
	}

	if tailModel.MultilineFlushInterval != nil {
		parent.InsertPairs("multiline_flush_interval", fmt.Sprint(*tailModel.MultilineFlushInterval))
	}

	if tailModel.PosFile != "" {
		parent.InsertPairs("pos_file", fmt.Sprint(tailModel.PosFile))
	}

	if tailModel.PosFileCompactionInterval != nil {
		parent.InsertPairs("pos_file_compaction_interval", fmt.Sprint(*tailModel.PosFileCompactionInterval))
	}

	if tailModel.PathKey != "" {
		parent.InsertPairs("path_key", fmt.Sprint(tailModel.PathKey))
	}

	if tailModel.RotateWait != nil {
		parent.InsertPairs("rotate_wait", fmt.Sprint(*tailModel.RotateWait))
	}

	if tailModel.EnableWatchTimer != nil {
		parent.InsertPairs("enable_watch_timer", fmt.Sprint(*tailModel.EnableWatchTimer))
	}

	if tailModel.EnableStatWatcher != nil {
		parent.InsertPairs("enable_stat_watcher", fmt.Sprint(*tailModel.EnableStatWatcher))
	}

	if tailModel.OpenOnEveryUpdate != nil {
		parent.InsertPairs("open_on_every_update", fmt.Sprint(*tailModel.OpenOnEveryUpdate))
	}

	if tailModel.EmitUnmatchedLines != nil {
		parent.InsertPairs("emit_unmatched_lines", fmt.Sprint(*tailModel.EmitUnmatchedLines))
	}

	if tailModel.IgnoreRepatedPermissionError != nil {
		parent.InsertPairs("ignore_repeated_permission_error", fmt.Sprint(*tailModel.IgnoreRepatedPermissionError))
	}

	return parent
}

func (i *Input) forwardPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	forwardModel := i.Forward
	childs := make([]*params.PluginStore, 0)
	if forwardModel.Transport != nil {
		child, _ := forwardModel.Transport.Params(loader)
		childs = append(childs, child)
	}

	if forwardModel.Security != nil {
		child, _ := forwardModel.Security.Params(loader)
		childs = append(childs, child)
	}

	if forwardModel.Client != nil {
		child, _ := forwardModel.Client.Params(loader)
		childs = append(childs, child)
	}

	if forwardModel.User != nil {
		if forwardModel.User.Username != nil && forwardModel.User.Password != nil {
			child, _ := forwardModel.User.Params(loader)
			childs = append(childs, child)
		}
	}

	parent.InsertChilds(childs...)

	if forwardModel.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*forwardModel.Port))
	}

	if forwardModel.Bind != nil {
		parent.InsertPairs("bind", fmt.Sprint(*forwardModel.Bind))
	}

	if forwardModel.Tag != nil {
		parent.InsertPairs("tag", fmt.Sprint(*forwardModel.Tag))
	}

	if forwardModel.AddTagPrefix != nil {
		parent.InsertPairs("add_tag_prefix", fmt.Sprint(*forwardModel.AddTagPrefix))
	}

	if forwardModel.LingerTimeout != nil {
		parent.InsertPairs("linger_timeout", fmt.Sprint(*forwardModel.LingerTimeout))
	}

	if forwardModel.ResolveHostname != nil {
		parent.InsertPairs("resolve_hostname", fmt.Sprint(*forwardModel.ResolveHostname))
	}

	if forwardModel.DenyKeepalive != nil {
		parent.InsertPairs("deny_keepalive", fmt.Sprint(*forwardModel.DenyKeepalive))
	}

	if forwardModel.SendKeepalivePacket != nil {
		parent.InsertPairs("send_keepalive_packet", fmt.Sprint(*forwardModel.SendKeepalivePacket))
	}

	if forwardModel.ChunkSizeLimit != nil {
		parent.InsertPairs("chunk_size_limit", fmt.Sprint(*forwardModel.ChunkSizeLimit))
	}

	if forwardModel.ChunkSizeWarnLimit != nil {
		parent.InsertPairs("chunk_size_warn_limit", fmt.Sprint(*forwardModel.ChunkSizeWarnLimit))
	}

	if forwardModel.SkipInvalidEvent != nil {
		parent.InsertPairs("skip_invalid_event", fmt.Sprint(*forwardModel.SkipInvalidEvent))
	}

	if forwardModel.SourceAddressKey != nil {
		parent.InsertPairs("source_address_key", fmt.Sprint(*forwardModel.SourceAddressKey))
	}

	return parent
}

func (i *Input) httpPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	httpModel := i.Http
	childs := make([]*params.PluginStore, 0)

	if httpModel.Transport != nil {
		child, _ := httpModel.Transport.Params(loader)
		childs = append(childs, child)
	}

	if httpModel.Parse != nil {
		child, _ := httpModel.Parse.Params(loader)
		childs = append(childs, child)
	}

	parent.InsertChilds(childs...)

	if httpModel.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*httpModel.Port))
	}

	if httpModel.Bind != nil {
		parent.InsertPairs("bind", fmt.Sprint(*httpModel.Bind))
	}

	if httpModel.BodySizeLimit != nil {
		parent.InsertPairs("body_size_limit", fmt.Sprint(*httpModel.BodySizeLimit))
	}

	if httpModel.KeepLiveTimeout != nil {
		parent.InsertPairs("keepalive_timeout", fmt.Sprint(*httpModel.KeepLiveTimeout))
	}

	if httpModel.AddHttpHeaders != nil {
		parent.InsertPairs("add_http_headers", fmt.Sprint(*httpModel.AddHttpHeaders))
	}

	if httpModel.AddRemoteAddr != nil {
		parent.InsertPairs("add_remote_addr", fmt.Sprint(*httpModel.AddRemoteAddr))
	}

	if httpModel.CorsAllowOrigins != nil {
		parent.InsertPairs("cors_allow_origins", fmt.Sprint(*httpModel.CorsAllowOrigins))
	}

	if httpModel.CorsAllowCredentials != nil {
		parent.InsertPairs("cors_allow_credentials", fmt.Sprint(*httpModel.CorsAllowCredentials))
	}

	if httpModel.RespondsWithEmptyImg != nil {
		parent.InsertPairs("responds_with_empty_img", fmt.Sprint(*httpModel.RespondsWithEmptyImg))
	}

	return parent
}

func (i *Input) samplePlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	sampleModel := i.Sample
	if sampleModel.Tag != nil {
		parent.InsertPairs("tag", fmt.Sprint(*sampleModel.Tag))
	}
	if sampleModel.Rate != nil {
		parent.InsertPairs("rate", fmt.Sprint(*sampleModel.Rate))
	}
	if sampleModel.Size != nil {
		parent.InsertPairs("size", fmt.Sprint(*sampleModel.Size))
	}
	if sampleModel.AutoIncrementKey != nil {
		parent.InsertPairs("auto_increment_key", fmt.Sprint(*sampleModel.AutoIncrementKey))
	}
	if sampleModel.Sample != nil {
		parent.InsertPairs("sample", fmt.Sprint(*sampleModel.Sample))
	}
	return parent
}

func (i *Input) monitorAgentPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	monitorAgentModel := i.MonitorAgent
	if monitorAgentModel.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*monitorAgentModel.Port))
	}
	if monitorAgentModel.Bind != nil {
		parent.InsertPairs("bind", fmt.Sprint(*monitorAgentModel.Bind))
	}
	if monitorAgentModel.Tag != nil {
		parent.InsertPairs("tag", fmt.Sprint(*monitorAgentModel.Tag))
	}
	if monitorAgentModel.EmitInterval != nil {
		parent.InsertPairs("emit_interval", fmt.Sprint(*monitorAgentModel.EmitInterval))
	}
	if monitorAgentModel.IncludeConfig != nil {
		parent.InsertPairs("include_config", fmt.Sprint(*monitorAgentModel.IncludeConfig))
	}
	if monitorAgentModel.IncludeRetry != nil {
		parent.InsertPairs("include_retry", fmt.Sprint(*monitorAgentModel.IncludeRetry))
	}
	return parent
}

var _ plugins.Plugin = &Input{}
