package client

import (
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/utils"
	"github.com/AlexandreMarcq/gozimbra/test"
	"github.com/google/go-cmp/cmp"
)

func TestParseAttributes(t *testing.T) {
	tests := map[string]struct {
		body       []byte
		attributes []string
		shouldErr  bool
		want       utils.AttrsMap
	}{
		"parsing attributes": {body: []byte(`<TestRequest><a n="ATTR1">VALUE1</a><a n="ATTR2">VALUE2</a></TestRequest>`), attributes: []string{"ATTR1", "ATTR2"}, shouldErr: false, want: utils.AttrsMap{"ATTR1": "VALUE1", "ATTR2": "VALUE2"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldErr {
				_, err := parseAttributes(tc.body, tc.attributes)

				test.AssertError(t, err)
			} else {
				got, err := parseAttributes(tc.body, tc.attributes)

				test.AssertNoError(t, err)

				diff := cmp.Diff(tc.want, got)
				if diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
