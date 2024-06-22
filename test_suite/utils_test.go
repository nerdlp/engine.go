package testsuite

import (
	"nerdlp/engine.go"
	"testing"
)

// TODO: need to write test for this
func TestComputePath(t *testing.T) {
	tests := map[string]struct {
		input          string
		handleResponse func(t *testing.T, out string)
	}{
		"no_leading_slash": {
			input: "engine.io",
			handleResponse: func(t *testing.T, out string) {
				if out != "/engine.io/" {
					t.Fatalf("expected \\engine.io\\, but got %s", out)
				}
			},
		},
		"leading_slashes": {
			input: "///engine.io",
			handleResponse: func(t *testing.T, out string) {
				if out != "/engine.io/" {
					t.Fatalf("expected \\engine.io\\, but got %s", out)
				}
			},
		},
		"trailing_slashes": {
			input: "engine.io////",
			handleResponse: func(t *testing.T, out string) {
				if out != "/engine.io/" {
					t.Fatalf("expected \\engine.io\\, but got %s", out)
				}
			},
		},
		"leading_slashes_and_trailing_slashes": {
			input: "///engine.io//",
			handleResponse: func(t *testing.T, out string) {
				if out != "/engine.io/" {
					t.Fatalf("expected \\engine.io\\, but got %s", out)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := engine.ComputePath(tt.input)
			tt.handleResponse(t, output)
		})
	}
}
