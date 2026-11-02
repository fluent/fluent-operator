package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

// The inline code is written on a single line, so a "--" comment left in it
// would swallow everything after it, including the final "end".
func TestFilter_Lua_InlineCode(t *testing.T) {
	cases := []struct {
		name string
		code string
		want string
	}{
		{
			name: "no comments",
			code: "function cb(tag, timestamp, record)\n" +
				"    record[\"cluster\"] = os.getenv(\"CLUSTER\") or \"unknown\"\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record["cluster"] = os.getenv("CLUSTER") or "unknown"; return 1, timestamp, record; end `,
		},
		{
			name: "comment on its own line",
			code: "function cb(tag, timestamp, record)\n" +
				"    -- add cluster metadata\n" +
				"    record[\"cluster\"] = os.getenv(\"CLUSTER\") or \"unknown\"\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record["cluster"] = os.getenv("CLUSTER") or "unknown"; return 1, timestamp, record; end `,
		},
		{
			name: "comment after code",
			code: "function cb(tag, timestamp, record)\n" +
				"    record[\"a\"] = 1 -- the answer\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record["a"] = 1; return 1, timestamp, record; end `,
		},
		{
			name: "dashes inside a double quoted string",
			code: "function cb(tag, timestamp, record)\n" +
				"    record[\"x\"] = \"a -- b\"\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record["x"] = "a -- b"; return 1, timestamp, record; end `,
		},
		{
			name: "dashes inside a single quoted string",
			code: "function cb(tag, timestamp, record)\n" +
				"    record['x'] = 'a -- b' -- note\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record['x'] = 'a -- b'; return 1, timestamp, record; end `,
		},
		{
			name: "dashes inside a long string",
			code: "function cb(tag, timestamp, record)\n" +
				"    record[\"x\"] = [[a -- b]]\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) record["x"] = [[a -- b]]; return 1, timestamp, record; end `,
		},
		{
			name: "block comment on one line is left as it was",
			code: "function cb(tag, timestamp, record)\n" +
				"    --[[ add cluster metadata ]]\n" +
				"    record[\"a\"] = 1\n" +
				"    return 1, timestamp, record\n" +
				"end\n",
			want: `function cb(tag, timestamp, record) --[[ add cluster metadata ]] record["a"] = 1; return 1, timestamp, record; end `,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			sl := plugins.NewSecretLoader(nil, "test namespace")
			l := Lua{Call: "cb", Code: tc.code}

			expected := params.NewKVs()
			expected.Insert("code", tc.want)
			expected.Insert("call", "cb")

			kvs, err := l.Params(sl)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(kvs).To(Equal(expected))
		})
	}
}
