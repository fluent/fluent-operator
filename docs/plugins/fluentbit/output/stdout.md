# Stdout

The stdout output plugin allows to print to the standard output the data received through the input plugin.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| format | Specify the data format to be printed. Supported formats are msgpack json, json_lines and json_stream. | string |
| jsonDateKey | Specify the name of the date field in output. | string |
| jsonDateFormat | Specify the format of the date. Supported formats are double,  iso8601 (eg: 2018-05-30T09:39:52.000681Z) and epoch. | string |
