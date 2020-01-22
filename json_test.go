package semver

import (
	"encoding/json"
	"testing"
)

func Test_JSON(t *testing.T) {
	tt := []struct {
		name string
		v    Version
	}{
		{
			name: "default",
			v:    Version{1, 2, 3},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.v)
			if err != nil {
				t.Errorf("marshal: %v", err)
			}
			var actual Version
			if err := json.Unmarshal(data, &actual); err != nil {
				t.Errorf("unmarshal: %v", err)
			}
			if !actual.Eq(tc.v) {
				t.Errorf("not equal")
			}
		})
	}
}
