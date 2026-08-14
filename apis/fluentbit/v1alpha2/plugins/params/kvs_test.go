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
			got, err := kvs.YamlString(tt.args.depth)
			if err != nil {
				t.Fatalf("YamlString() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("YamlString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestKVs_String_RejectsNewlineInjection reproduces GHSA-2j8x-46rv-qmpq for Fluent
// Bit's classic renderer: a newline embedded in a value would let the value break
// onto a new line and inject a top-level [OUTPUT] section. The classic format has
// no escaping, so such values must be rejected.
func TestKVs_String_RejectsNewlineInjection(t *testing.T) {
	kvs := NewKVs()
	kvs.Insert("Name", "http")
	kvs.Insert("Host", "http://legit.com\n\n[OUTPUT]\n    Name    file\n    Path    /tmp/exfil")

	if _, err := kvs.String(); err == nil {
		t.Fatal("expected String() to reject a value containing a newline, got nil error")
	}
}

func TestKVs_YamlString_RejectsNewlineInjection(t *testing.T) {
	kvs := NewKVs()
	kvs.Insert("Name", "http")
	kvs.Insert("Host", "http://legit.com\ninjected: true")

	if _, err := kvs.YamlString(1); err == nil {
		t.Fatal("expected YamlString() to reject a value containing a newline, got nil error")
	}
}

// TestKVs_String_AllowsMultilineContent ensures the raw Content passthrough (used
// for CustomPlugin bodies, which are legitimately multi-line) is not rejected.
func TestKVs_String_AllowsMultilineContent(t *testing.T) {
	kvs := NewKVs()
	kvs.Content = "    Name    cpu\n    Tag    my_cpu\n"

	got, err := kvs.String()
	if err != nil {
		t.Fatalf("String() unexpected error for Content passthrough: %v", err)
	}
	if got != kvs.Content {
		t.Errorf("String() = %q, want %q", got, kvs.Content)
	}
}
