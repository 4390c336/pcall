package pcall

import (
	"testing"

	"github.com/coredns/caddy"
)

func Test_setup(t *testing.T) {

	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			"run linux.example.org",
			`pcall {
				run ./test/resolver
			}`,
			false,
		},
		{
			"run bad op",
			`pcall {
				runn ./test/resolver
			}`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := caddy.NewTestController("dns", tt.config)
			if err := setup(ctr); (err != nil) != tt.wantErr {
				t.Errorf("setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
