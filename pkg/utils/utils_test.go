package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestContainString(t *testing.T) {
	type args struct {
		item  string
		array []string
	}
	tests := []struct {
		name   string
		args   args
		wanted bool
	}{{
		name: "target array not contain",
		args: args{
			item:  "a",
			array: []string{"c", "b", "d"},
		},
		wanted: false,
	}, {
		name: "target array not contain one item",
		args: args{
			item:  "a",
			array: []string{"c", "b", "d", "a"},
		},
		wanted: true,
	}, {
		name: "target array not contain more items",
		args: args{
			item:  "a",
			array: []string{"a", "c", "a", "a"},
		},
		wanted: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainString(tt.args.array, tt.args.item); !reflect.DeepEqual(got, tt.wanted) {
				t.Errorf("The return of ContainString() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	type args struct {
		item  string
		array []string
	}
	tests := []struct {
		name   string
		args   args
		wanted []string
	}{{
		name: "target array contains one item, and removed",
		args: args{
			item:  "a",
			array: []string{"b", "a", "c", "d"},
		},
		wanted: []string{"b", "c", "d"},
	}, {
		name: "target array contains more items, and removed",
		args: args{
			item:  "a",
			array: []string{"a", "b", "a", "c", "d", "a"},
		},
		wanted: []string{"b", "c", "d"},
	}, {
		name: "target array not contain the item",
		args: args{
			item:  "f",
			array: []string{"a", "c", "d", "e"},
		},
		wanted: []string{"a", "c", "d", "e"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveString(tt.args.array, tt.args.item); !reflect.DeepEqual(got, tt.wanted) {
				t.Errorf("The return of RemoveString() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestConcatString(t *testing.T) {
	type args struct {
		sep   string
		array []string
	}
	tests := []struct {
		name   string
		args   args
		wanted string
	}{{
		name: "the slice is nil",
		args: args{
			sep:   "|",
			array: nil,
		},
		wanted: "",
	}, {
		name: "the slice is empty",
		args: args{
			sep:   "|",
			array: []string{},
		},
		wanted: "",
	}, {
		name: "concat a string with |",
		args: args{
			sep:   "|",
			array: []string{"b", "a", "c", "d"},
		},
		wanted: "b|a|c|d",
	}, {
		name: "concat a string with a space",
		args: args{
			sep:   " ",
			array: []string{"b", "a", "c", "d"},
		},
		wanted: "b a c d",
	}, {
		name: "concat a string",
		args: args{
			sep:   "",
			array: []string{"l", "o", "v", "e"},
		},
		wanted: "love",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatString(tt.args.array, tt.args.sep); !reflect.DeepEqual(got, tt.wanted) {
				t.Errorf("The return of ConcatString() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestHashCode(t *testing.T) {
	fmt.Println(HashCode("hello"))
}
