# Tail

The in_tail Input plugin allows Fluentd to read events from the tail of text files. Its behavior is similar to the tail -F command.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | The tag of the event. | string |
| path | The path(s) to read. Multiple paths can be specified, separated by comma ','. | string |
| pathTimezone | This parameter is for strftime formatted path like /path/to/%Y/%m/%d/. | string |
| excludePath | The paths excluded from the watcher list. | []string |
| followInodes | Avoid to read rotated files duplicately. You should set true when you use * or strftime format in path. | *bool |
| refreshInterval | The interval to refresh the list of watch files. This is used when the path includes *. | *uint32 |
| limitRecentlyModified | Limits the watching files that the modification time is within the specified time range when using * in path. | *uint32 |
| skipRefreshOnStartup | Skips the refresh of the watch list on startup. This reduces the startup time when * is used in path. | *bool |
| readFromHead | Starts to read the logs from the head of the file or the last read position recorded in pos_file, not tail. | *bool |
| encoding | Specifies the encoding of reading lines. By default, in_tail emits string value as ASCII-8BIT encoding. If encoding is specified, in_tail changes string to encoding. If encoding and fromEncoding both are specified, in_tail tries to encode string from fromEncoding to encoding. | string |
| fromEncoding | Specifies the encoding of reading lines. By default, in_tail emits string value as ASCII-8BIT encoding. If encoding is specified, in_tail changes string to encoding. If encoding and fromEncoding both are specified, in_tail tries to encode string from fromEncoding to encoding. | string |
| readLinesLimit | The number of lines to read with each I/O operation. | *int32 |
| readBytesLimitPerSecond | The number of reading bytes per second to read with I/O operation. This value should be equal or greater than 8192. | *int32 |
| maxLineSize | The maximum length of a line. Longer lines than it will be just skipped. | *int32 |
| multilineFlushInterval | The interval of flushing the buffer for multiline format. | *uint32 |
| posFile | (recommended) Fluentd will record the position it last read from this file. pos_file handles multiple positions in one file so no need to have multiple pos_file parameters per source. Don't share pos_file between in_tail configurations. It causes unexpected behavior e.g. corrupt pos_file content. | string |
| posFileCompactionInterval | The interval of doing compaction of pos file. | *uint32 |
| parse |  | *common.Parse |
| pathKey | Adds the watching file path to the path_key field. | string |
| rotateWait | in_tail actually does a bit more than tail -F itself. When rotating a file, some data may still need to be written to the old file as opposed to the new one. in_tail takes care of this by keeping a reference to the old file (even after it has been rotated) for some time before transitioning completely to the new file. This helps prevent data designated for the old file from getting lost. By default, this time interval is 5 seconds. The rotate_wait parameter accepts a single integer representing the number of seconds you want this time interval to be. | *uint32 |
| enableWatchTimer | Enables the additional watch timer. Setting this parameter to false will significantly reduce CPU and I/O consumption when tailing a large number of files on systems with inotify support. The default is true which results in an additional 1 second timer being used. | *bool |
| enableStatWatcher | Enables the additional inotify-based watcher. Setting this parameter to false will disable the inotify events and use only timer watcher for file tailing. This option is mainly for avoiding the stuck issue with inotify. | *bool |
| openOnEveryUpdate | Opens and closes the file on every update instead of leaving it open until it gets rotated. | *bool |
| emitUnmatchedLines | Emits unmatched lines when <parse> format is not matched for incoming logs. | *bool |
| ignoreRepeatedPermissionError | If you have to exclude the non-permission files from the watch list, set this parameter to true. It suppresses the repeated permission error logs. | *bool |
| group | The in_tail plugin can assign each log file to a group, based on user defined rules. The limit parameter controls the total number of lines collected for a group within a rate_period time interval. | *Group |
# Group




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| pattern | Specifies the regular expression for extracting metadata (namespace, podname) from log file path. Default value of the pattern regexp extracts information about namespace, podname, docker_id, container of the log (K8s specific). | string |
| ratePeriod | Time period in which the group line limit is applied. in_tail resets the counter after every rate_period interval. | *int32 |
| rule | Grouping rules for log files. | *Rule |
# Rule




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | match parameter is used to check if a file belongs to a particular group based on hash keys (named captures from pattern) and hash values (regexp in string) | map[string]string |
| limit | Maximum number of lines allowed from a group in rate_period time interval. The default value of -1 doesn't throttle log files of that group. | *int32 |
