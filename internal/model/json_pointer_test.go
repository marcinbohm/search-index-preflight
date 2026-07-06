package model

import "testing"

func TestAppendJSONPointer(t *testing.T) {
	tests := []struct {
		name   string
		base   string
		tokens []string
		want   string
	}{
		{
			name:   "empty base",
			base:   "",
			tokens: []string{"dynamic"},
			want:   "/dynamic",
		},
		{
			name:   "non-empty base",
			base:   "/mappings",
			tokens: []string{"dynamic"},
			want:   "/mappings/dynamic",
		},
		{
			name:   "multiple tokens",
			base:   "/template/mappings",
			tokens: []string{"properties", "status"},
			want:   "/template/mappings/properties/status",
		},
		{
			name:   "slash escaping",
			base:   "",
			tokens: []string{"service/name"},
			want:   "/service~1name",
		},
		{
			name:   "tilde escaping",
			base:   "",
			tokens: []string{"field~name"},
			want:   "/field~0name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendJSONPointer(tt.base, tt.tokens...); got != tt.want {
				t.Fatalf("AppendJSONPointer(%q, %v) = %q, want %q", tt.base, tt.tokens, got, tt.want)
			}
		})
	}
}
