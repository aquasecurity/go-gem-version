package gem

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     []string
	}{
		{
			name: "happy path",
			versions: []string{
				"1.1.1",
				"1.0",
				"1.2",
				"2",
				"0.7.1",
			},
			want: []string{
				"0.7.1",
				"1.0",
				"1.1.1",
				"1.2",
				"2",
			},
		},
		{
			name: "pre-release",
			versions: []string{
				"1.0.0.b",
				"1.0.0",
				"1.0.0.pre",
				"1.0.0-alpha.2",
				"1.0.0.1",
				"1.0.0.a.1",
				"1.0.0-alpha.11",
				"1.0.0.a",
				"1.1",
			},
			want: []string{
				"1.0.0.a",
				"1.0.0.a.1",
				"1.0.0.b",
				"1.0.0.pre.alpha.2",
				"1.0.0.pre.alpha.11",
				"1.0.0.pre",
				"1.0.0",
				"1.0.0.1",
				"1.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versions := make([]Version, len(tt.versions))
			for i, raw := range tt.versions {
				v, err := NewVersion(raw)
				require.NoError(t, err)
				versions[i] = v
			}

			sort.Sort(Collection(versions))

			got := make([]string, len(versions))
			for i, v := range versions {
				got[i] = v.String()
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
