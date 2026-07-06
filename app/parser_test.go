package main

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "double quotes preserve internal spaces",
			input: `echo "hello   world"`,
			want:  []string{"echo", "hello   world"},
		},
		{
			name:  "single quotes preserve internal spaces plus trailing token",
			input: `echo 'a b'   c`,
			want:  []string{"echo", "a b", "c"},
		},
		{
			name:  "bare command with no arguments",
			input: `echo`,
			want:  []string{"echo"},
		},
		{
			name:    "unterminated double quote returns error",
			input:   `echo "hello`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("tokenize(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("tokenize(%q) unexpected error: %v", tt.input, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
