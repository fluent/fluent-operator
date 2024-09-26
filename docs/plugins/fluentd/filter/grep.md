# Grep

Grep defines various parameters for the grep plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | []*[Regexp](#regexp) |
| exclude |  | []*[Exclude](#exclude) |
| and |  | []*[And](#and) |
| or |  | []*[Or](#or) |
# Regexp

Regexp defines the parameters for the regexp plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key |  | *string |
| pattern |  | *string |
# Exclude

Exclude defines the parameters for the exclude plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key |  | *string |
| pattern |  | *string |
# And

And defines the parameters for the \"and\" plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | *[Regexp](#regexp) |
| exclude |  | *[Exclude](#exclude) |
# Or

Or defines the parameters for the \"or\" plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regexp |  | *[Regexp](#regexp) |
| exclude |  | *[Exclude](#exclude) |
