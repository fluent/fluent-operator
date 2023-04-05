# File

The file output plugin allows to write the data received through the input plugin to file. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/file**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| path | Absolute directory path to store files. If not set, Fluent Bit will write the files on it's own positioned directory. | string |
| file | Set file name to store the records. If not set, the file name will be the tag associated with the records. | string |
| format | The format of the file content. See also Format section. Default: out_file. | string |
| delimiter | The character to separate each pair. Applicable only if format is csv or ltsv. | string |
| labelDelimiter | The character to separate each pair. Applicable only if format is ltsv. | string |
| template | The format string. Applicable only if format is template. | string |
