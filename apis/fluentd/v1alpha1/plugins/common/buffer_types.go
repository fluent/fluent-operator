package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/params"
)

// +kubebuilder:object:generate:=true
// BufferCommon defines common parameters for the buffer plugin
type BufferCommon struct {
	// The @id parameter specifies a unique name for the configuration.
	Id *string `json:"id,omitempty"`
	// The @type parameter specifies the type of the plugin.
	// +kubebuilder:validation:Enum:=file;memory;file_single
	// +kubebuilder:validation:Required
	Type *string `json:"type"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
}

// Buffer defines various parameters for the buffer Plugin
type Buffer struct {
	BufferCommon `json:",inline,omitempty"`
	// The file buffer plugin
	*FileBuffer `json:",inline,omitempty"`
	// The file_single buffer plugin
	*FileSingleBuffer `json:",inline,omitempty"`
	// The time section of buffer plugin
	Time *Time `json:",inline,omitempty"`
	// Output plugin will flush chunks per specified time (enabled when time is specified in chunk keys)
	TimeKey *string `json:"timekey,omitempty"`
	// Output plugin will write chunks after timekey_wait seconds later after timekey expiration
	TimeKeyWait *string `json:"timekeyWait,omitempty"`

	// The path where buffer chunks are stored. This field would make no effect in memory buffer plugin.
	Path *string `json:"path,omitempty"`
	// The output plugins group events into chunks.
	// Chunk keys, specified as the argument of <buffer> section, control how to group events into chunks.
	// If tag is empty, which means blank Chunk Keys.
	// Tag also supports Nested Field, combination of Chunk Keys, placeholders, etc.
	// See https://docs.fluentd.org/configuration/buffer-section.
	Tag string `json:"tag,omitempty"`
	// Buffer parameters
	// The max size of each chunks: events will be written into chunks until the size of chunks become this size
	// Default: 8MB (memory) / 256MB (file)
	// +kubebuilder:validation:Pattern:="^\\d+(KB|MB|GB|TB)$"
	ChunkLimitSize *string `json:"chunkLimitSize,omitempty"`
	// The max number of events that each chunks can store in it.
	// +kubebuilder:validation:Pattern:="^\\d+(KB|MB|GB|TB)$"
	ChunkLimitRecords *string `json:"chunkLimitRecords,omitempty"`
	// The size limitation of this buffer plugin instance
	// Default: 512MB (memory) / 64GB (file)
	// +kubebuilder:validation:Pattern:="^\\d+(KB|MB|GB|TB)$"
	TotalLimitSize *string `json:"totalLimitSize,omitempty"`
	// The queue length limitation of this buffer plugin instance. Default: 0.95
	// +kubebuilder:validation:Pattern:="^\\d+.?\\d+$"
	QueueLimitLength *string `json:"queueLimitLength,omitempty"`
	// Limit the number of queued chunks. Default: 1
	// If a smaller flush_interval is set, e.g. 1s,
	// there are lots of small queued chunks in the buffer.
	// With file buffer, it may consume a lot of fd resources when output destination has a problem.
	// This parameter mitigates such situations.
	// +kubebuilder:validation:Minimum:=1
	QueuedChunksLimitSize *int16 `json:"queuedChunksLimitSize,omitempty"`
	// Fluentd will decompress these compressed chunks automatically before passing them to the output plugin
	// If gzip is set, Fluentd compresses data records before writing to buffer chunks.
	// Default:text.
	// +kubebuilder:validation:Enum:=text;gzip
	Compress *string `json:"compress,omitempty"`
	// Flush parameters
	// This specifies whether to flush/write all buffer chunks on shutdown or not.
	FlushAtShutdown *bool `json:"flushAtShutdown,omitempty"`
	// FlushMode defines the flush mode:
	// lazy: flushes/writes chunks once per timekey
	// interval: flushes/writes chunks per specified time via flush_interval
	// immediate: flushes/writes chunks immediately after events are appended into chunks
	// default: equals to lazy if time is specified as chunk key, interval otherwise
	// +kubebuilder:validation:Enum:=default;lazy;interval;immediate
	FlushMode *string `json:"flushMode,omitempty"`
	// FlushInterval defines the flush interval
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	FlushInterval *string `json:"flushInterval,omitempty"`
	// The number of threads to flush/write chunks in parallel
	// +kubebuilder:validation:Pattern:="^\\d+$"
	FlushThreadCount *string `json:"flushThreadCount,omitempty"`
	// The timeout (seconds) until output plugin decides if the async write operation has failed. Default is 60s
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	DelayedCommitTimeout *string `json:"delayedCommitTimeout,omitempty"`
	// OverflowAtction defines the output plugin behave when its buffer queue is full.
	// +kubebuilder:validation:Enum:throw_exception,block,drop_oldest_chunk
	// Default: throw_exception
	OverflowAction *string `json:"overflowAction,omitempty"`
	// Retry parameters
	// The maximum time (seconds) to retry to flush again the failed chunks, until the plugin discards the buffer chunks
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	RetryTimeout *string `json:"retryTimeout,omitempty"`
	// If true, plugin will ignore retry_timeout and retry_max_times options and retry flushing forever.
	RetryForever *bool `json:"retryForever,omitempty"`
	// The maximum number of times to retry to flush the failed chunks. Default: none
	RetryMaxTimes *int16 `json:"retryMaxTimes,omitempty"`
	// The ratio of retry_timeout to switch to use the secondary while failing.
	// +kubebuilder:validation:Pattern:="^\\d+.?\\d+$"
	RetrySecondaryThreshold *string `json:"retrySecondaryThreshold,omitempty"`
	// Output plugin will retry periodically with fixed intervals.
	// +kubebuilder:validation:Enum:exponential_backoff,periodic
	RetryType *string `json:"retryType,omitempty"`
	// Wait in seconds before the next retry to flush or constant factor of exponential backoff
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	RetryWait *string `json:"retryWait,omitempty"`
	// The base number of exponential backoff for retries.
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?$"
	RetryExponentialBackoffBase *string `json:"retryExponentialBackoffBase,omitempty"`
	// The maximum interval (seconds) for exponential backoff between retries while failing
	// +kubebuilder:validation:Pattern:="^\\d+(\\.[0-9]{0,2})?(s|m|h|d)?$"
	RetryMaxInterval *string `json:"retryMaxInterval,omitempty"`
	// If true, the output plugin will retry after randomized interval not to do burst retries
	RetryRandomize *bool `json:"retryRandomize,omitempty"`
	// Instead of storing unrecoverable chunks in the backup directory, just discard them. This option is new in Fluentd v1.2.6.
	DisableChunkBackup *bool `json:"disableChunkBackup,omitempty"`
}

// The file buffer plugin provides a persistent buffer implementation. It uses files to store buffer chunks on disk.
type FileBuffer struct {
	// Changes the suffix of the buffer file.
	PathSuffix *string `json:"pathSuffix,omitempty"`
}

// The file_single buffer plugin is similar to file_file but it does not have the metadata file.
// See https://docs.fluentd.org/buffer/file_single#limitation
type FileSingleBuffer struct {
	// Calculates the number of records, chunk size, during chunk resume.
	CalcNumRecords *string `json:"calcNumRecords,omitempty"`
	// ChunkFormat specifies the chunk format for calc_num_records.
	// +kubebuilder:validation:Enum:=msgpack;text;auto
	ChunkFormat *string `json:"chunkFormat,omitempty"`
}

// BufferSection defines the common parameters for buffer sections
type BufferSection struct {
	// buffer section
	Buffer *Buffer `json:"buffer,omitempty"`
	// format section
	Format *Format `json:"format,omitempty"`
	// inject section
	Inject *Inject `json:"inject,omitempty"`
}

func (b *Buffer) Name() string {
	return "buffer"
}

func (b *Buffer) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(b.Name())
	params.InsertPairs(ps, "@id", b.Id)
	if b.Type != nil {
		ps.InsertType(fmt.Sprint(*b.Type))
	}
	params.InsertPairs(ps, "@log_level", b.LogLevel)

	if b.FileBuffer != nil && b.PathSuffix != nil {
		ps.InsertPairs("path_suffix", *b.PathSuffix)
	}

	if b.FileSingleBuffer != nil {
		params.InsertPairs(ps, "calc_num_records", b.CalcNumRecords)
		params.InsertPairs(ps, "chunk_format", b.ChunkFormat)
	}

	if b.Path != nil {
		if strings.HasPrefix(*b.Path, "/buffers") {
			ps.InsertPairs("path", *b.Path)
		} else {
			targetPaths := []string{"/buffers"}
			paths := strings.Split(*b.Path, "/")
			targetPaths = append(targetPaths, paths...)
			ps.InsertPairs("path", filepath.Join(targetPaths...))
		}
	} else {
		if *b.Type != "memory" {
			ps.InsertPairs("path", params.DefaultBufferPath)
		}
	}

	params.InsertPairs(ps, "timekey", b.TimeKey)
	params.InsertPairs(ps, "timekey_wait", b.TimeKeyWait)

	ps.InsertPairs("tag", b.Tag)

	params.InsertPairs(ps, "chunk_limit_size", b.ChunkLimitSize)
	params.InsertPairs(ps, "chunk_limit_records", b.ChunkLimitRecords)
	params.InsertPairs(ps, "total_limit_size", b.TotalLimitSize)
	params.InsertPairs(ps, "queue_limit_length", b.QueueLimitLength)
	params.InsertPairs(ps, "queued_chunks_limit_size", b.QueuedChunksLimitSize)
	params.InsertPairs(ps, "compress", b.Compress)

	if b.FlushAtShutdown != nil && *b.FlushAtShutdown {
		ps.InsertPairs("flush_at_shutdown", fmt.Sprint(*b.FlushAtShutdown))
	}

	params.InsertPairs(ps, "flush_mode", b.FlushMode)
	params.InsertPairs(ps, "flush_interval", b.FlushInterval)
	params.InsertPairs(ps, "flush_thread_count", b.FlushThreadCount)
	params.InsertPairs(ps, "delayed_commit_timeout", b.DelayedCommitTimeout)
	params.InsertPairs(ps, "overflow_action", b.OverflowAction)
	params.InsertPairs(ps, "retry_timeout", b.RetryTimeout)
	params.InsertPairs(ps, "retry_secondary_threshold", b.RetrySecondaryThreshold)
	params.InsertPairs(ps, "retry_type", b.RetryType)
	params.InsertPairs(ps, "retry_max_times", b.RetryMaxTimes)
	params.InsertPairs(ps, "retry_forever", b.RetryForever)
	params.InsertPairs(ps, "retry_wait", b.RetryWait)
	params.InsertPairs(ps, "retry_exponential_backoff_base", b.RetryExponentialBackoffBase)
	params.InsertPairs(ps, "retry_max_interval", b.RetryMaxInterval)
	params.InsertPairs(ps, "retry_randomize", b.RetryRandomize)
	params.InsertPairs(ps, "disable_chunk_backup", b.DisableChunkBackup)

	if b.Time != nil {
		params.InsertPairs(ps, "time_type", b.Time.TimeType)
		params.InsertPairs(ps, "time_format", b.Time.TimeFormat)
		params.InsertPairs(ps, "localtime", b.Time.Localtime)
		params.InsertPairs(ps, "utc", b.Time.UTC)
		params.InsertPairs(ps, "timezone", b.Time.Timezone)
		params.InsertPairs(ps, "time_format_fallbacks", b.Time.TimeFormatFallbacks)
	}

	return ps, nil
}
