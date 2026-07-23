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
		{
			name:  "unquoted escaped spaces stay literal within one word",
			input: `echo three\ \ \ spaces`,
			want:  []string{"echo", "three   spaces"},
		},
		{
			name:  "unquoted escaped space mixed with unescaped spaces",
			input: "echo before\\   after",
			want:  []string{"echo", "before ", "after"},
		},
		{
			name:  "unquoted backslash removes special meaning of ordinary letter",
			input: `echo test\nexample`,
			want:  []string{"echo", "testnexample"},
		},
		{
			name:  "unquoted double backslash produces one literal backslash",
			input: `echo hello\\world`,
			want:  []string{"echo", "hello\\world"},
		},
		{
			name:  "unquoted backslash escapes single quote characters",
			input: `echo \'hello\'`,
			want:  []string{"echo", "'hello'"},
		},
		{
			name:  "single quotes keep backslash literal",
			input: `echo 'a\nb'`,
			want:  []string{"echo", `a\nb`},
		},
		{
			name:  "double quotes escape dollar sign",
			input: `echo "a\$b"`,
			want:  []string{"echo", "a$b"},
		},
		{
			name:  "double quotes escape embedded double quote",
			input: `echo "a\"b"`,
			want:  []string{"echo", `a"b`},
		},
		{
			name:  "double quotes escape backslash",
			input: `echo "a\\b"`,
			want:  []string{"echo", `a\b`},
		},
		{
			name:  "double quotes keep backslash literal before unrelated char",
			input: `echo "a\ab"`,
			want:  []string{"echo", `a\ab`},
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
