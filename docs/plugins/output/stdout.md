# Stdout

The stdout output plugin allows to print to the standard output the data received through the input plugin.


| Field | Description | Scheme | Default |
| ----- | ----------- | ------ | ----- |
| format | Specify the data format to be printed. Supported formats are msgpack json, json_lines and json_stream. | string | msgpack |
| jsonDateKey | Specify the name of the date field in output.To disable the time key just set the value to `false`. | string | date |
| jsonDateFormat | Specify the format of the date. Supported formats are double,  iso8601 (eg: 2018-05-30T09:39:52.000681Z) and epoch. | string | double |
