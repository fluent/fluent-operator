package filter

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	v1 "k8s.io/api/core/v1"
)

// +kubebuilder:object:generate:=true

// The Lua Filter allows you to modify the incoming records using custom Lua Scripts. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/lua**
type Lua struct {
	plugins.CommonParams `json:",inline"`
	// Path to the Lua script that will be used.
	Script v1.ConfigMapKeySelector `json:"script,omitempty"`
	// Lua function name that will be triggered to do filtering.
	// It's assumed that the function is declared inside the Script defined above.
	Call string `json:"call"`
	// Inline LUA code instead of loading from a path via script.
	Code string `json:"code,omitempty"`
	// If these keys are matched, the fields are converted to integer.
	// If more than one key, delimit by space.
	// Note that starting from Fluent Bit v1.6 integer data types are preserved
	// and not converted to double as in previous versions.
	TypeIntKey []string `json:"typeIntKey,omitempty"`
	// If these keys are matched, the fields are handled as array. If more than
	// one key, delimit by space. It is useful the array can be empty.
	TypeArrayKey []string `json:"typeArrayKey,omitempty"`
	// If enabled, Lua script will be executed in protected mode.
	// It prevents to crash when invalid Lua script is executed. Default is true.
	ProtectedMode *bool `json:"protectedMode,omitempty"`
	// If enabled, null will be converted to flb_null in Lua.
	// This helps prevent removing key/value since nil is a special value to remove key/value from map in Lua. Default value: false.
	EnableFlbNull *bool `json:"enable_flb_null,omitempty"`
	// By default when the Lua script is invoked, the record timestamp is passed as a
	// Floating number which might lead to loss precision when the data is converted back.
	// If you desire timestamp precision enabling this option will pass the timestamp as
	// a Lua table with keys sec for seconds since epoch and nsec for nanoseconds.
	TimeAsTable bool `json:"timeAsTable,omitempty"`
}

func (l *Lua) Name() string {
	return "lua"
}

func (l *Lua) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := l.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}

	if l.Code != "" {
		var singleLineLua = ""
		lineTrim := ""
		re := regexp.MustCompile(`^function |^if |^for |^else|^elseif |^end|--[[]+`)
		for line := range strings.SplitSeq(strings.TrimSuffix(l.Code, "\n"), "\n") {
			lineTrim = strings.TrimSpace(line)
			if lineTrim != "" {
				operator := re.MatchString(lineTrim)
				if operator {
					singleLineLua = singleLineLua + lineTrim + " "
				} else {
					singleLineLua = singleLineLua + lineTrim + "; "
				}
			}
		}
		kvs.Insert("code", singleLineLua)
	}

	if l.Script.Key != "" {
		kvs.Insert("script", "/fluent-bit/config/"+l.Script.Key)
	}

	kvs.Insert("call", l.Call)

	if len(l.TypeIntKey) > 0 {
		kvs.Insert("type_int_key", strings.Join(l.TypeIntKey, " "))
	}

	if len(l.TypeArrayKey) > 0 {
		kvs.Insert("type_array_key", strings.Join(l.TypeArrayKey, " "))
	}

	if l.ProtectedMode != nil {
		kvs.Insert("protected_mode", strconv.FormatBool(*l.ProtectedMode))
	}

	if l.EnableFlbNull != nil {
		kvs.Insert("enable_flb_null", strconv.FormatBool(*l.EnableFlbNull))
	}

	if l.TimeAsTable {
		kvs.Insert("time_as_table", "true")
	}

	return kvs, nil
}

func (l *Lua) MakeNamespaced(ns string) {
	l.Script.Key = fmt.Sprintf("%x-%s", md5.Sum([]byte(ns)), l.Script.Key)
}
