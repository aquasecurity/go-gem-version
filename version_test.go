package gem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		version string
		wantErr bool
	}{
		{"1.2.3", false},
		{"1", false},
		{"1.2.beta", false},
		{"1.21.beta", false},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"1.2.0-x.Y.0+metadata", true},
		{"1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", true},
		{"1.2.0.4-x.Y.0+metadata-width-hypen", true},
		{"1.2.0-X-1.2.0+metadata~dist", true},
		{"1.2.3.4-rc1-with-hypen", false},
		{"foo1.2.3", true},
		{"1.7rc2", false},
		{"v1.7rc2", true},
		{"1.0-", true},

		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L120-L124
		{"", false},
		{" ", false},
		{"   ", false},

		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L83-L89
		{"1.0", false},
		{"1.0 ", false},
		{" 1.0 ", false},
		{"\n1.0", false},
		{"1.0\n", false},
		{"\n1.0\n", false},

		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L91-L109
		{"junk", true},
		{"1.0\\n2.0", true},
		{"1..2", true},
		{"1.2\\ 3.4", true},
		{"1.2\\ 3.4", true},
		{"2.3422222.222.222222222.22222.ads0as.dasd0.ddd2222.2.qd3e.", true},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			_, err := NewVersion(tt.version)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want int
	}{
		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L149-L167
		{"1.0", "1.0.0", 0},
		{"1.0", "1.0.a", 1},
		{"1.8.2", "0.0.0", 1},
		{"1.8.2", "1.8.2.a", 1},
		{"1.8.2.b", "1.8.2.a", 1},
		{"1.8.2.a", "1.8.2", -1},
		{"1.8.2.a10", "1.8.2.a9", 1},
		{"", "0", 0},

		{"0.beta.1", "0.0.beta.1", 0},
		{"0.0.beta", "0.0.beta.1", -1},
		{"0.0.beta", "0.beta.1", -1},

		{"5.a", "5.0.0.rc2", -1},
		{"5.x", "5.0.0.rc2", 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.v1, tt.v2), func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			require.NoError(t, err, tt.v1)

			v2, err := NewVersion(tt.v2)
			require.NoError(t, err, tt.v2)

			assert.Equal(t, tt.want, v1.Compare(v2))
		})
	}
}

func TestVersion_Equal(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L69-L73
		{"1.2", "1.2", true},
		{"1.2", "1.2.0", true},
		{"1.2", "1.3", false},
		{"1.2.b1", "1.2.b.1", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.v1, tt.v2), func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			require.NoError(t, err, tt.v1)

			v2, err := NewVersion(tt.v2)
			require.NoError(t, err, tt.v2)

			assert.Equal(t, tt.want, v1.Equal(v2))
		})
	}
}
func TestVersion_LessThan(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L196-L203
		{"1.0.0-alpha", "1.0.0-alpha.1", true},
		{"1.0.0-alpha.1", "1.0.0-beta.2", true},
		{"1.0.0-beta.2", "1.0.0-beta.11", true},
		{"1.0.0-beta.11", "1.0.0-rc.1", true},
		{"1.0.0-rc1", "1.0.0", true},
		{"1.0.0-1", "1", true},
		{"1", "1.0.0-1", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.v1, tt.v2), func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			require.NoError(t, err, tt.v1)

			v2, err := NewVersion(tt.v2)
			require.NoError(t, err, tt.v2)

			assert.Equal(t, tt.want, v1.LessThan(v2))
		})
	}
}
func TestVersion_GreaterThan(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		// https://github.com/rubygems/rubygems/blob/6914b4ec439ae1e7079b3c08576cb3fbce68aa69/test/rubygems/test_gem_version.rb#L149-L167
		{"1.0", "1.0.0", false},
		{"1.0", "1.0.a", true},
		{"1.8.2", "0.0.0", true},
		{"1.8.2", "1.8.2.a", true},
		{"1.8.2.b", "1.8.2.a", true},
		{"1.8.2.a", "1.8.2", false},
		{"1.8.2.a10", "1.8.2.a9", true},
		{"", "0", false},

		{"0.beta.1", "0.0.beta.1", false},
		{"0.0.beta", "0.0.beta.1", false},
		{"0.0.beta", "0.beta.1", false},

		{"5.a", "5.0.0.rc2", false},
		{"5.x", "5.0.0.rc2", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.v1, tt.v2), func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			require.NoError(t, err, tt.v1)

			v2, err := NewVersion(tt.v2)
			require.NoError(t, err, tt.v2)

			assert.Equal(t, tt.want, v1.GreaterThan(v2))
		})
	}
}

func TestVersion_Release(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"5.3.1", "5.3.1"},
		{"1.2.3.4", "1.2.3.4"},
		{"5.3.1.b.2", "5.3.1"},
		{"1", "1"},
		{"1.a", "1"},
		{"1.2", "1.2"},
		{"1.2.a.3", "1.2"},

		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_version.rb#L143-L148
		{"1.2.0.a", "1.2.0"},
		{"1.1.rc10", "1.1"},
		{"1.9.3.alpha.5", "1.9.3"},
		{"1.9.3", "1.9.3"},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v, err := NewVersion(tt.version)
			require.NoError(t, err)

			want, err := NewVersion(tt.want)
			require.NoError(t, err)

			got := v.Release()
			assert.Equal(t, want.canonicalSegments(), got.canonicalSegments())
		})
	}
}

func TestVersion_Bump(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_version.rb#L12-L30
		{"5.2.4", "5.3"},
		{"5.2.4.a", "5.3"},
		{"5.2.4.a10", "5.3"},
		{"5.0.0", "5.1"},
		{"5.0", "6"},
		{"5", "6"},
		{"1.2.3.4", "1.2.4"},
		{"5.3.1.b.2", "5.4"},
		{"1.a", "2"},
		{"1.2", "2"},
		{"1.2.a.3", "2"},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v, err := NewVersion(tt.version)
			require.NoError(t, err)

			want, err := NewVersion(tt.want)
			require.NoError(t, err)

			got := v.Bump()
			assert.Equal(t, want.canonicalSegments(), got.canonicalSegments())
		})
	}
}
