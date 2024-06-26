package params

import "testing"

func TestKVs_YamlString(t *testing.T) {
	type fields struct {
		keys    []string
		values  []string
		Content string
	}
	type args struct {
		depth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				keys:   []string{"Daemon", "Flush", "Grace"},
				values: []string{"false", "5", "30"},
			},
			args: args{
				depth: 1,
			},
			want: "  daemon: false\n  flush: 5\n  grace: 30\n",
		},
		{
			name: "Test 2",
			fields: fields{
				keys:   []string{"Remove_key", "Remove_key"},
				values: []string{"stream", "time"},
			},
			args: args{
				depth: 1,
			},
			want: "  remove_key:\n    - stream\n    - time\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs := &KVs{
				keys:    tt.fields.keys,
				values:  tt.fields.values,
				Content: tt.fields.Content,
			}
			if got := kvs.YamlString(tt.args.depth); got != tt.want {
				t.Errorf("YamlString() = %v, want %v", got, tt.want)
			}
		})
	}
}
