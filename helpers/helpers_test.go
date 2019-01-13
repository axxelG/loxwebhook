package helpers

import (
	"testing"
)

func Test_IsStringInSlice(t *testing.T) {
	list := []string{
		"a",
		"b",
		"c",
	}
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "found",
			str:  "b",
			want: true,
		},
		{
			name: "notFound",
			str:  "z",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStringInSlice(tt.str, list); got != tt.want {
				t.Errorf("isStringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapStringKeyFromStringValue(t *testing.T) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	tests := []struct {
		name    string
		str     string
		wantKey string
		wantOk  bool
	}{
		{
			name:    "found",
			str:     "value2",
			wantKey: "key2",
			wantOk:  true,
		},
		{
			name:    "found",
			str:     "value99",
			wantKey: "",
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotOk := GetMapStringKeyFromStringValue(tt.str, m)
			if gotKey != tt.wantKey {
				t.Errorf("GetMapStringKeyFromStringValue() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetMapStringKeyFromStringValue() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
