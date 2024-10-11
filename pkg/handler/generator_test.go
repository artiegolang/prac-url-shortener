package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateShortID(t *testing.T) {
	type args struct {
		longURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{longURL: "https://google.com"},
			want: GenerateShortID("https://google.com"),
		},
		{
			name: "2",
			args: args{longURL: "https://.com"},
			want: GenerateShortID("https://.com"),
		},
		{
			name: "empty",
			args: args{longURL: ""},
			want: GenerateShortID(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, GenerateShortID(tt.args.longURL), tt.want)
		})
	}
}
