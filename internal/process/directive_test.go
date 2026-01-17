package process_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tahardi/pluckmd/internal/process"
)

func TestNewDirective(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{
			name:    "valid - standard type",
			line:    `<!-- pluck("go", "type", "Verifier", "tee/verifier.go", 0, 0) -->`,
			wantErr: false,
		},
		{
			name:    "valid - function with negative indices",
			line:    `<!-- pluck("go", "function", "Verify", "tee/verifier.go", -1, -1) -->`,
			wantErr: false,
		},
		{
			name:    "valid - messy spacing",
			line:    `<!--pluck(   "go" ,  "type"  ,"Name"  ,  "path" , 0 , 0 )-->`,
			wantErr: false,
		},
		{
			name:    "invalid - wrong number of fields",
			line:    `<!-- pluck("type", "Name", "path", 0) -->`,
			wantErr: true,
		},
		{
			name:    "invalid - non-integer index",
			line:    `<!-- pluck("go", "type", "Name", "path", 0, "end") -->`,
			wantErr: true,
		},
		{
			name:    "invalid - unknown lang",
			line:    `<!-- pluck("not-a-lang", "type", "Name", "path", 0, 0) -->`,
			wantErr: true,
		},
		{
			name:    "invalid - unknown kind",
			line:    `<!-- pluck("go", "not-a-kind", "Name", "path", 0, 0) -->`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := process.NewDirective(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
