# BufferCommon

BufferCommon defines common parameters for the buffer plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| type | The @type parameter specifies the type of the plugin. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Buffer

Buffer defines various parameters for the buffer Plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| timekey | Output plugin will flush chunks per specified time (enabled when time is specified in chunk keys) | *string |
| timekeyWait | Output plugin will write chunks after timekey_wait seconds later after timekey expiration | *string |
| path | The path where buffer chunks are stored. This field would make no effect in memory buffer plugin. | *string |
| tag | The output plugins group events into chunks. Chunk keys, specified as the argument of <buffer> section, control how to group events into chunks. If tag is empty, which means blank Chunk Keys. Tag also supports Nested Field, combination of Chunk Keys, placeholders, etc. See https://docs.fluentd.org/configuration/buffer-section. | string |
| chunkLimitSize | Buffer parameters The max size of each chunks: events will be written into chunks until the size of chunks become this size Default: 8MB (memory) / 256MB (file) | *string |
| chunkLimitRecords | The max number of events that each chunks can store in it. | *string |
| totalLimitSize | The size limitation of this buffer plugin instance Default: 512MB (memory) / 64GB (file) | *string |
| queueLimitLength | The queue length limitation of this buffer plugin instance. Default: 0.95 | *string |
| queuedChunksLimitSize | Limit the number of queued chunks. Default: 1 If a smaller flush_interval is set, e.g. 1s, there are lots of small queued chunks in the buffer. With file buffer, it may consume a lot of fd resources when output destination has a problem. This parameter mitigates such situations. | *int16 |
| compress | Fluentd will decompress these compressed chunks automatically before passing them to the output plugin If gzip is set, Fluentd compresses data records before writing to buffer chunks. Default:text. | *string |
| flushAtShutdown | Flush parameters This specifies whether to flush/write all buffer chunks on shutdown or not. | *bool |
| flushMode | FlushMode defines the flush mode: lazy: flushes/writes chunks once per timekey interval: flushes/writes chunks per specified time via flush_interval immediate: flushes/writes chunks immediately after events are appended into chunks default: equals to lazy if time is specified as chunk key, interval otherwise | *string |
| flushInterval | FlushInterval defines the flush interval | *string |
| flushThreadCount | The number of threads to flush/write chunks in parallel | *string |
| delayedCommitTimeout | The timeout (seconds) until output plugin decides if the async write operation has failed. Default is 60s | *string |
| overflowAction | OverflowAtction defines the output plugin behave when its buffer queue is full. Default: throw_exception | *string |
| retryTimeout | Retry parameters The maximum time (seconds) to retry to flush again the failed chunks, until the plugin discards the buffer chunks | *string |
| retryForever | If true, plugin will ignore retry_timeout and retry_max_times options and retry flushing forever. | *bool |
| retryMaxTimes | The maximum number of times to retry to flush the failed chunks. Default: none | *int16 |
| retrySecondaryThreshold | The ratio of retry_timeout to switch to use the secondary while failing. | *string |
| retryType | Output plugin will retry periodically with fixed intervals. | *string |
| retryWait | Wait in seconds before the next retry to flush or constant factor of exponential backoff | *string |
| retryExponentialBackoffBase | The base number of exponential backoff for retries. | *string |
| retryMaxInterval | The maximum interval (seconds) for exponential backoff between retries while failing | *string |
| retryRandomize | If true, the output plugin will retry after randomized interval not to do burst retries | *bool |
| disableChunkBackup | Instead of storing unrecoverable chunks in the backup directory, just discard them. This option is new in Fluentd v1.2.6. | *bool |
# FileBuffer

The file buffer plugin provides a persistent buffer implementation. It uses files to store buffer chunks on disk.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| pathSuffix | Changes the suffix of the buffer file. | *string |
# FileSingleBuffer

The file_single buffer plugin is similar to file_file but it does not have the metadata file. See https://docs.fluentd.org/buffer/file_single#limitation


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| calcNumRecords | Calculates the number of records, chunk size, during chunk resume. | *string |
| chunkFormat | ChunkFormat specifies the chunk format for calc_num_records. | *string |
# BufferSection

BufferSection defines the common parameters for buffer sections


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| buffer | buffer section | *Buffer |
| format | format section | *Format |
| inject | inject section | *Inject |
