# Systemd

The Systemd input plugin allows to collect log messages from the Journald daemon on Linux environments. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/systemd**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| path | Optional path to the Systemd journal directory, if not set, the plugin will use default paths to read local-only logs. | string |
| db | Specify the database file to keep track of monitored files and offsets. | string |
| dbSync | Set a default synchronization (I/O) method. values: Extra, Full, Normal, Off. This flag affects how the internal SQLite engine do synchronization to disk, for more details about each option please refer to this section. note: this option was introduced on Fluent Bit v1.4.6. | string |
| tag | The tag is used to route messages but on Systemd plugin there is an extra functionality: if the tag includes a star/wildcard, it will be expanded with the Systemd Unit file (e.g: host.* => host.UNIT_NAME). | string |
| maxFields | Set a maximum number of fields (keys) allowed per record. | int |
| maxEntries | When Fluent Bit starts, the Journal might have a high number of logs in the queue. In order to avoid delays and reduce memory usage, this option allows to specify the maximum number of log entries that can be processed per round. Once the limit is reached, Fluent Bit will continue processing the remaining log entries once Journald performs the notification. | int |
| systemdFilter | Allows to perform a query over logs that contains a specific Journald key/value pairs, e.g: _SYSTEMD_UNIT=UNIT. The Systemd_Filter option can be specified multiple times in the input section to apply multiple filters as required. | []string |
| systemdFilterType | Define the filter type when Systemd_Filter is specified multiple times. Allowed values are And and Or. With And a record is matched only when all of the Systemd_Filter have a match. With Or a record is matched when any of the Systemd_Filter has a match. | string |
| readFromTail | Start reading new entries. Skip entries already stored in Journald. | string |
| stripUnderscores | Remove the leading underscore of the Journald field (key). For example the Journald field _PID becomes the key PID. | string |
| storageType | Specify the buffering mechanism to use. It can be memory or filesystem | string |
| pauseOnChunksOverlimit | Specifies if the input plugin should be paused (stop ingesting new data) when the storage.max_chunks_up value is reached. | string |
