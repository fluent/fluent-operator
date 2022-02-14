# Grep

Grep defines all supported types for filter_grep plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | []*Regexp |
| exclude |  | []*Exclude |
| and |  | []*And |
| or |  | []*Or |
# Regexp

Regexp defines the parameters for regexp plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key |  | *string |
| pattern |  | *string |
# Exclude

Exclude defines the parameters for exclude plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key |  | *string |
| pattern |  | *string |
# And

And defines the parameters for and plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | *Regexp |
| exclude |  | *Exclude |
# Or

Or defines the parameters for or plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | *Regexp |
| exclude |  | *Exclude |
