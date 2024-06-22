package testsuite

import (
	"nerdlp/engine.go"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		input struct {
			a int
			b int
		}
		handleResponse func(t *testing.T, output int)
	}{
		"add_1_and_3_return_3": {
			input: struct {
				a int
				b int
			}{
				a: 1,
				b: 2,
			},
			handleResponse: func(t *testing.T, output int) {
				if output != 4 {
					t.Fatalf("expected 3, but got %v", output)
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := engine.Add(tt.input.a, tt.input.b)
			tt.handleResponse(t, output)
		})
	}
}
