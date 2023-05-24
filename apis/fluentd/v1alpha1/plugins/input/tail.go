package input

import (
	"encoding/json"
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/common"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/params"
)

// The in_tail Input plugin allows Fluentd to read events from the tail of text files. Its behavior is similar to the tail -F command.
type Tail struct {
	// +kubebuilder:validation:Required
	// The tag of the event.
	Tag string `json:"tag"`
	// +kubebuilder:validation:Required
	// The path(s) to read. Multiple paths can be specified, separated by comma ','.
	Path string `json:"path"`
	// This parameter is for strftime formatted path like /path/to/%Y/%m/%d/.
	PathTimezone string `json:"pathTimezone,omitempty"`
	// The paths excluded from the watcher list.
	ExcludePath []string `json:"excludePath,omitempty"`
	// Avoid to read rotated files duplicately. You should set true when you use * or strftime format in path.
	FollowInodes *bool `json:"followInodes,omitempty"`
	// The interval to refresh the list of watch files. This is used when the path includes *.
	RefreshInterval *uint32 `json:"refreshInterval,omitempty"`
	// Limits the watching files that the modification time is within the specified time range when using * in path.
	LimitRecentlyModified *uint32 `json:"limitRecentlyModified,omitempty"`
	// Skips the refresh of the watch list on startup. This reduces the startup time when * is used in path.
	SkipRefreshOnStartup *bool `json:"skipRefreshOnStartup,omitempty"`
	// Starts to read the logs from the head of the file or the last read position recorded in pos_file, not tail.
	ReadFromHead *bool `json:"readFromHead,omitempty"`
	// Specifies the encoding of reading lines. By default, in_tail emits string value as ASCII-8BIT encoding.
	// If encoding is specified, in_tail changes string to encoding.
	// If encoding and fromEncoding both are specified, in_tail tries to encode string from fromEncoding to encoding.
	Encoding string `json:"encoding,omitempty"`
	// Specifies the encoding of reading lines. By default, in_tail emits string value as ASCII-8BIT encoding.
	// If encoding is specified, in_tail changes string to encoding.
	// If encoding and fromEncoding both are specified, in_tail tries to encode string from fromEncoding to encoding.
	FromEncoding string `json:"fromEncoding,omitempty"`
	// The number of lines to read with each I/O operation.
	ReadLinesLimit *int32 `json:"readLinesLimit,omitempty"`
	// The number of reading bytes per second to read with I/O operation. This value should be equal or greater than 8192.
	ReadBytesLimitPerSecond *int32 `json:"readBytesLimitPerSecond,omitempty"`
	// The maximum length of a line. Longer lines than it will be just skipped.
	MaxLineSize *int32 `json:"maxLineSize,omitempty"`
	// The interval of flushing the buffer for multiline format.
	MultilineFlushInterval *uint32 `json:"multilineFlushInterval,omitempty"`
	// (recommended) Fluentd will record the position it last read from this file.
	// pos_file handles multiple positions in one file so no need to have multiple pos_file parameters per source.
	// Don't share pos_file between in_tail configurations. It causes unexpected behavior e.g. corrupt pos_file content.
	PosFile string `json:"posFile,omitempty"`
	// The interval of doing compaction of pos file.
	PosFileCompactionInterval *uint32 `json:"posFileCompactionInterval,omitempty"`
	// +kubebuilder:validation:Required
	Parse *common.Parse `json:"parse"`
	// Adds the watching file path to the path_key field.
	PathKey string `json:"pathKey,omitempty"`
	// in_tail actually does a bit more than tail -F itself. When rotating a file, some data may still need to be written to the old file as opposed to the new one.
	// in_tail takes care of this by keeping a reference to the old file (even after it has been rotated) for some time before transitioning completely to the new file.
	// This helps prevent data designated for the old file from getting lost. By default, this time interval is 5 seconds.
	// The rotate_wait parameter accepts a single integer representing the number of seconds you want this time interval to be.
	RotateWait *uint32 `json:"rotateWait,omitempty"`
	// Enables the additional watch timer. Setting this parameter to false will significantly reduce CPU and I/O consumption when tailing a large number of files on systems with inotify support.
	// The default is true which results in an additional 1 second timer being used.
	EnableWatchTimer *bool `json:"enableWatchTimer,omitempty"`
	// Enables the additional inotify-based watcher. Setting this parameter to false will disable the inotify events and use only timer watcher for file tailing.
	// This option is mainly for avoiding the stuck issue with inotify.
	EnableStatWatcher *bool `json:"enableStatWatcher,omitempty"`
	// Opens and closes the file on every update instead of leaving it open until it gets rotated.
	OpenOnEveryUpdate *bool `json:"openOnEveryUpdate,omitempty"`
	// Emits unmatched lines when <parse> format is not matched for incoming logs.
	EmitUnmatchedLines *bool `json:"emitUnmatchedLines,omitempty"`
	// If you have to exclude the non-permission files from the watch list, set this parameter to true. It suppresses the repeated permission error logs.
	IgnoreRepatedPermissionError *bool `json:"ignoreRepeatedPermissionError,omitempty"`
	// The in_tail plugin can assign each log file to a group, based on user defined rules.
	// The limit parameter controls the total number of lines collected for a group within a rate_period time interval.
	Group *Group `json:"group,omitempty"`
}

type Group struct {
	// Specifies the regular expression for extracting metadata (namespace, podname) from log file path.
	// Default value of the pattern regexp extracts information about namespace, podname, docker_id, container of the log (K8s specific).
	Pattern string `json:"pattern,omitempty"`
	// Time period in which the group line limit is applied. in_tail resets the counter after every rate_period interval.
	RatePeriod *int32 `json:"ratePeriod,omitempty"`
	// Grouping rules for log files.
	// +kubebuilder:validation:Required
	Rule *Rule `json:"rule"`
}

type Rule struct {
	// match parameter is used to check if a file belongs to a particular group based on hash keys (named captures from pattern) and hash values (regexp in string)
	Match map[string]string `json:"match,omitempty"`
	// Maximum number of lines allowed from a group in rate_period time interval. The default value of -1 doesn't throttle log files of that group.
	Limit *int32 `json:"limit,omitempty"`
}

func (g *Group) Name() string {
	return "group"
}

func (g *Group) Params(loader plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(g.Name())

	if g.Pattern != "" {
		ps.InsertPairs("pattern", fmt.Sprint(g.Pattern))
	}

	if g.RatePeriod != nil {
		ps.InsertPairs("rate_period", fmt.Sprint(*g.RatePeriod))
	}

	if g.Rule != nil {
		subchild, _ := g.Rule.Params(loader)
		ps.InsertChilds(subchild)
	}
	return ps, nil
}

func (r *Rule) Name() string {
	return "rule"
}

func (r *Rule) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(r.Name())

	if len(r.Match) > 0 {
		matches, _ := json.Marshal(r.Match)
		ps.InsertPairs("match", fmt.Sprint(string(matches)))
	}

	if r.Limit != nil {
		ps.InsertPairs("limit", fmt.Sprint(*r.Limit))
	}

	return ps, nil
}
