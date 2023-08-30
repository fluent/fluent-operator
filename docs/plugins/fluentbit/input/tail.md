# Tail

The Tail input plugin allows to monitor one or several text files. <br /> It has a similar behavior like tail -f shell command. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/tail**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| bufferChunkSize | Set the initial buffer size to read files data. This value is used too to increase buffer size. The value must be according to the Unit Size specification. | string |
| bufferMaxSize | Set the limit of the buffer size per monitored file. When a buffer needs to be increased (e.g: very long lines), this value is used to restrict how much the memory buffer can grow. If reading a file exceed this limit, the file is removed from the monitored file list The value must be according to the Unit Size specification. | string |
| path | Pattern specifying a specific log files or multiple ones through the use of common wildcards. | string |
| pathKey | If enabled, it appends the name of the monitored file as part of the record. The value assigned becomes the key in the map. | string |
| excludePath | Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria, e.g: exclude_path=*.gz,*.zip | string |
| readFromHead | For new discovered files on start (without a database offset/position), read the content from the head of the file, not tail. | *bool |
| refreshIntervalSeconds | The interval of refreshing the list of watched files in seconds. | *int64 |
| rotateWaitSeconds | Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed. | *int64 |
| ignoredOlder | Ignores records which are older than this time in seconds. Supports m,h,d (minutes, hours, days) syntax. Default behavior is to read all records from specified files. Only available when a Parser is specificied and it can parse the time of a record. | string |
| skipLongLines | When a monitored file reach it buffer capacity due to a very long line (Buffer_Max_Size), the default behavior is to stop monitoring that file. Skip_Long_Lines alter that behavior and instruct Fluent Bit to skip long lines and continue processing other lines that fits into the buffer size. | *bool |
| db | Specify the database file to keep track of monitored files and offsets. | string |
| dbSync | Set a default synchronization (I/O) method. Values: Extra, Full, Normal, Off. | string |
| memBufLimit | Set a limit of memory that Tail plugin can use when appending data to the Engine. If the limit is reach, it will be paused; when the data is flushed it resumes. | string |
| parser | Specify the name of a parser to interpret the entry as a structured message. | string |
| key | When a message is unstructured (no parser applied), it's appended as a string under the key name log. This option allows to define an alternative name for that key. | string |
| tag | Set a tag (with regex-extract fields) that will be placed on lines read. E.g. kube.<namespace_name>.<pod_name>.<container_name> | string |
| tagRegex | Set a regex to exctract fields from the file | string |
| multiline | If enabled, the plugin will try to discover multiline messages and use the proper parsers to compose the outgoing messages. Note that when this option is enabled the Parser option is not used. | *bool |
| multilineFlushSeconds | Wait period time in seconds to process queued multiline messages | *int64 |
| parserFirstline | Name of the parser that matchs the beginning of a multiline message. Note that the regular expression defined in the parser must include a group name (named capture) | string |
| parserN | Optional-extra parser to interpret and structure multiline entries. This option can be used to define multiple parsers. | []string |
| dockerMode | If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above. This mode cannot be used at the same time as Multiline. | *bool |
| dockerModeFlushSeconds | Wait period time in seconds to flush queued unfinished split lines. | *int64 |
| dockerModeParser | Specify an optional parser for the first line of the docker multiline mode. The parser name to be specified must be registered in the parsers.conf file. | string |
| disableInotifyWatcher | DisableInotifyWatcher will disable inotify and use the file stat watcher instead. | *bool |
| multilineParser | This will help to reassembly multiline messages originally split by Docker or CRI Specify one or Multiline Parser definition to apply to the content. | string |
| storageType | Specify the buffering mechanism to use. It can be memory or filesystem | string |
| pauseOnChunksOverlimit | Specifies if the input plugin should be paused (stop ingesting new data) when the storage.max_chunks_up value is reached. | string |
