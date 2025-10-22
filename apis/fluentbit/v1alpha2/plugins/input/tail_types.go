package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Tail input plugin allows to monitor one or several text files. <br />
// It has a similar behavior like tail -f shell command. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/tail**
type Tail struct {
	// Set the initial buffer size to read files data.
	// This value is used too to increase buffer size.
	// The value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferChunkSize,omitempty"`
	// Set the limit of the buffer size per monitored file.
	// When a buffer needs to be increased (e.g: very long lines),
	// this value is used to restrict how much the memory buffer can grow.
	// If reading a file exceed this limit, the file is removed from the monitored file list
	// The value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// Pattern specifying a specific log files or multiple ones through the use of common wildcards.
	Path string `json:"path,omitempty"`
	// If enabled, it appends the name of the monitored file as part of the record.
	// The value assigned becomes the key in the map.
	PathKey string `json:"pathKey,omitempty"`
	// Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria,
	// e.g: exclude_path=*.gz,*.zip
	ExcludePath string `json:"excludePath,omitempty"`
	// If enabled, Fluent Bit appends the offset of the current monitored file as part of the record.
	// The value assigned becomes the key in the map
	OffsetKey string `json:"offsetKey,omitempty"`
	// For new discovered files on start (without a database offset/position),
	// read the content from the head of the file, not tail.
	ReadFromHead *bool `json:"readFromHead,omitempty"`
	// The interval of refreshing the list of watched files in seconds.
	RefreshIntervalSeconds *int64 `json:"refreshIntervalSeconds,omitempty"`
	// Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed.
	RotateWaitSeconds *int64 `json:"rotateWaitSeconds,omitempty"`
	// Ignores records which are older than this time in seconds.
	// Supports m,h,d (minutes, hours, days) syntax.
	// Default behavior is to read all records from specified files.
	// Only available when a Parser is specificied and it can parse the time of a record.
	// +kubebuilder:validation:Pattern:="^\\d+(m|h|d)?$"
	IgnoreOlder string `json:"ignoredOlder,omitempty"`
	// When a monitored file reach it buffer capacity due to a very long line (Buffer_Max_Size),
	// the default behavior is to stop monitoring that file.
	// Skip_Long_Lines alter that behavior and instruct Fluent Bit to skip long lines
	// and continue processing other lines that fits into the buffer size.
	SkipLongLines *bool `json:"skipLongLines,omitempty"`
	// Specify the database file to keep track of monitored files and offsets.
	DB string `json:"db,omitempty"`
	// Set a default synchronization (I/O) method. Values: Extra, Full, Normal, Off.
	// +kubebuilder:validation:Enum:=Extra;Full;Normal;Off
	DBSync string `json:"dbSync,omitempty"`
	// Specify that the database will be accessed only by Fluent Bit.
	DBLocking *bool `json:"dbLocking,omitempty"`
	// Set a limit of memory that Tail plugin can use when appending data to the Engine.
	// If the limit is reach, it will be paused; when the data is flushed it resumes.
	MemBufLimit string `json:"memBufLimit,omitempty"`
	// Specify the name of a parser to interpret the entry as a structured message.
	Parser string `json:"parser,omitempty"`
	// When a message is unstructured (no parser applied), it's appended as a string under the key name log.
	// This option allows to define an alternative name for that key.
	Key string `json:"key,omitempty"`
	// Set a tag (with regex-extract fields) that will be placed on lines read.
	// E.g. kube.<namespace_name>.<pod_name>.<container_name>
	Tag string `json:"tag,omitempty"`
	// Set a regex to exctract fields from the file
	TagRegex string `json:"tagRegex,omitempty"`
	// If enabled, the plugin will try to discover multiline messages
	// and use the proper parsers to compose the outgoing messages.
	// Note that when this option is enabled the Parser option is not used.
	Multiline *bool `json:"multiline,omitempty"`
	// Wait period time in seconds to process queued multiline messages
	MultilineFlushSeconds *int64 `json:"multilineFlushSeconds,omitempty"`
	// Name of the parser that matchs the beginning of a multiline message.
	// Note that the regular expression defined in the parser must include a group name (named capture)
	ParserFirstline string `json:"parserFirstline,omitempty"`
	// Optional-extra parser to interpret and structure multiline entries.
	// This option can be used to define multiple parsers.
	ParserN []string `json:"parserN,omitempty"`
	// If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above.
	// This mode cannot be used at the same time as Multiline.
	DockerMode *bool `json:"dockerMode,omitempty"`
	// Wait period time in seconds to flush queued unfinished split lines.
	DockerModeFlushSeconds *int64 `json:"dockerModeFlushSeconds,omitempty"`
	// Specify an optional parser for the first line of the docker multiline mode. The parser name to be specified must be registered in the parsers.conf file.
	DockerModeParser string `json:"dockerModeParser,omitempty"`
	// DisableInotifyWatcher will disable inotify and use the file stat watcher instead.
	DisableInotifyWatcher *bool `json:"disableInotifyWatcher,omitempty"`
	// This will help to reassembly multiline messages originally split by Docker or CRI
	// Specify one or Multiline Parser definition to apply to the content.
	MultilineParser string `json:"multilineParser,omitempty"`
	// Specify the buffering mechanism to use. It can be memory or filesystem
	// +kubebuilder:validation:Enum:=filesystem;memory
	StorageType string `json:"storageType,omitempty"`
	// Specifies if the input plugin should be paused (stop ingesting new data) when the storage.max_chunks_up value is reached.
	// +kubebuilder:validation:Enum:=on;off
	PauseOnChunksOverlimit string `json:"pauseOnChunksOverlimit,omitempty"`
	// Skips empty lines in the log file from any further processing or output.
	SkipEmptyLines *bool `json:"skipEmptyLines,omitempty"`
	// Threaded mechanism allows input plugin to run in a separate thread which helps to desaturate the main pipeline.
	Threaded *string `json:"threaded,omitempty"`
}

func (*Tail) Name() string {
	return "tail"
}

func (t *Tail) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Buffer_Chunk_Size", t.BufferChunkSize)
	plugins.InsertKVString(kvs, "Buffer_Max_Size", t.BufferMaxSize)
	plugins.InsertKVString(kvs, "Path", t.Path)
	plugins.InsertKVString(kvs, "Path_Key", t.PathKey)
	plugins.InsertKVString(kvs, "Exclude_Path", t.ExcludePath)
	plugins.InsertKVString(kvs, "Offset_Key", t.OffsetKey)
	plugins.InsertKVField(kvs, "Refresh_Interval", t.RefreshIntervalSeconds)
	plugins.InsertKVString(kvs, "Ignore_Older", t.IgnoreOlder)
	plugins.InsertKVField(kvs, "Skip_Long_Lines", t.SkipLongLines)
	plugins.InsertKVString(kvs, "DB", t.DB)
	plugins.InsertKVString(kvs, "DB.Sync", t.DBSync)
	plugins.InsertKVString(kvs, "Mem_Buf_Limit", t.MemBufLimit)
	plugins.InsertKVString(kvs, "Parser", t.Parser)
	plugins.InsertKVString(kvs, "Key", t.Key)
	plugins.InsertKVString(kvs, "Tag", t.Tag)
	plugins.InsertKVString(kvs, "Tag_Regex", t.TagRegex)
	plugins.InsertKVField(kvs, "Docker_Mode", t.DockerMode)
	plugins.InsertKVField(kvs, "Docker_Mode_Flush", t.DockerModeFlushSeconds)
	plugins.InsertKVString(kvs, "Docker_Mode_Parser", t.DockerModeParser)
	plugins.InsertKVString(kvs, "Parser_Firstline", t.ParserFirstline)
	plugins.InsertKVString(kvs, "multiline.parser", t.MultilineParser)
	plugins.InsertKVString(kvs, "storage.type", t.StorageType)
	plugins.InsertKVString(kvs, "storage.pause_on_chunks_overlimit", t.PauseOnChunksOverlimit)

	plugins.InsertKVField(kvs, "Read_from_Head", t.ReadFromHead)
	plugins.InsertKVField(kvs, "Rotate_Wait", t.RotateWaitSeconds)
	plugins.InsertKVField(kvs, "DB.locking", t.DBLocking)
	plugins.InsertKVField(kvs, "Multiline", t.Multiline)
	plugins.InsertKVField(kvs, "Multiline_Flush", t.MultilineFlushSeconds)
	plugins.InsertKVField(kvs, "Skip_Empty_Lines", t.SkipEmptyLines)

	for i, parser := range t.ParserN {
		kvs.Insert(fmt.Sprintf("Parser_%d", i+1), parser)
	}
	if t.DisableInotifyWatcher != nil {
		kvs.Insert("Inotify_Watcher", fmt.Sprint(!*t.DisableInotifyWatcher))
	}
	if t.Threaded != nil {
		kvs.Insert("Threaded", *t.Threaded)
	}

	return kvs, nil
}
