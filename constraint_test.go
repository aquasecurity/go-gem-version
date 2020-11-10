package gem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConstraints(t *testing.T) {
	tests := []struct {
		constraint string
		wantErr    bool
	}{
		{"> 1.0", false},
		{"> 1.0 || < foo", true},

		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L276-L282
		{">>> 1.3.5", true},
		{"> blah", true},
	}
	for _, tt := range tests {
		t.Run(tt.constraint, func(t *testing.T) {
			_, err := NewConstraints(tt.constraint)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVersion_Check(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		// Not equal
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L117-L127
		{"!= 1.2", "1.1", true},
		{"!= 1.2", "1.2", false},
		{"!= 1.2", "1.3", true},

		// Blank
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L129-L139
		{"1.2", "1.1", false},
		{"1.2", "1.2", true},
		{"1.2", "1.3", false},

		// Equal
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L141-L151
		{"= 1.2", "1.1", false},
		{"=1.2", "1.2", true},
		{"= 1.2", "1.3", false},

		// Equal: ==
		{"== 1.2", "1.1", false},
		{"==1.2", "1.2", true},
		{"== 1.2", "1.3", false},

		// Greater than
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L153-L163
		{"> 1.2", "1.1", false},
		{">1.2", "1.2", false},
		{"> 1.2", "1.3", true},

		// Greater than or equal
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L165-L175
		{">= 1.2", "1.1", false},
		{">=1.2", "1.2", true},
		{">= 1.2", "1.3", true},

		// List: comma separated
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L177-L187
		{"> 1.1, < 1.3", "1.1", false},
		{"> 1.1, <1.3", "1.2", true},
		{"> 1.1, < 1.3", "1.3", false},

		// List: space separated
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L177-L187
		{"> 1.1 < 1.3", "1.1", false},
		{"> 1.1	<1.3", "1.2", true},
		{"> 1.1 < 1.3", "1.3", false},

		// Less than
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L177-L187
		{"< 1.2", "1.1", true},
		{"<1.2", "1.2", false},
		{"< 1.2", "1.3", false},

		// Less than or equal
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L201-L211
		{"<= 1.2", "1.1", true},
		{"<=1.2", "1.2", true},
		{"<= 1.2", "1.3", false},

		// Pessimistic
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L213-L223
		{"~> 1.2", "1.1", false},
		{"~>1.2", "1.2", true},
		{"~> 1.2", "1.3", true},

		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L225-L231
		{"~> 0.0.1", "0.1.1", false},
		{"~>0.0.1", "0.0.2", true},
		{"~> 0.0.1", "0.0.1", true},

		// Good
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L233-L274
		{"= 0.2.33", "0.2.33", true},
		{"== 0.2.33", "0.2.33", true},
		{"> 0.2.33", "0.2.34", true},
		{"= 1.0", "1.0", true},
		{"= 1.0", "1.0.0", true},
		{"= 1.0.0", "1.0", true},
		{"1.0", "1.0", true},
		{"> 1.8.0", "1.8.2", true},
		{"> 1.111", "1.112", true},
		{"> 0.0.0", "0.2", true},
		{"> 0.0.0", "0.0.0.0.0.2", true},
		{"> 0.0.0.1", "0.0.1.0", true},
		{"> 9.3.2", "10.3.2", true},
		{"= 1.0", "1.0.0.0", true},
		{"!= 9.3.4", "10.3.2", true},
		{"> 9.3.2", "10.3.2", true},
		{">= 9.3.2", " 9.3.2", true},
		{">= 9.3.2", "9.3.2", true},
		{"= 0", "", true},
		{"< 0.1", "", true},
		{"< 0.1", "	", true},
		{" < 0.1", "", true},
		{"> 0.a", "	", true},
		{">	0.a", "", true},
		{"< 3.2.rc1", "3.1", true},

		{"> 3.2.0.rc1", "3.2.0", true},
		{"> 3.2.0.rc1", "3.2.0.rc2", true},

		{"< 3.0", "3.0.rc2", true},
		{"< 3.0.0", "3.0.rc2", true},
		{"< 3.0.1", "3.0.rc2", true},

		{"> 0", "3.0.rc2", true},

		{"~> 5.a", "5.0.0.rc2", true},
		{"~> 5.x", "5.0.0.rc2", false},

		{"~> 5.a", "5.0.0", true},
		{"~> 5.x", "5.0.0", true},

		// Boxed
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L294-L316
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L329-L342
		{"~> 1.4", "1.3", false},
		{"~> 1.4", "1.4", true},
		{"~> 1.4", "1.4.0", true},
		{"~> 1.4", "1.5", true},
		{"~> 1.4", "2.0", false},

		{"~> 1.4.4", "1.3", false},
		{"~> 1.4.4", "1.4", false},
		{"~> 1.4.4", "1.4.4", true},
		{"~> 1.4.4", "1.4.5", true},
		{"~> 1.4.4", "1.5", false},
		{"~> 1.4.4", "2.0", false},

		{"~> 1.0.0", "1.1.pre", false},
		{"~> 1.1", "1.1.pre", false},
		{"~> 1.0", "2.0.a", false},
		{"~> 2.0", "2.0.a", false},

		{"~> 1", "0.9", false},
		{"~> 1", "1.0", true},
		{"~> 1", "1.1", true},
		{"~> 1", "2.0", false},

		// Multiple
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L318-L327
		{">= 1.4, <= 1.6, != 1.5", "1.3", false},
		{">= 1.4  <= 1.6  != 1.5", "1.4", true},
		{">= 1.4  <= 1.6  != 1.5", "1.5", false},
		{">= 1.4, <= 1.6, != 1.5", "1.6", true},
		{">= 1.4, <= 1.6, != 1.5", "1.7", false},
		{">= 1.4, <= 1.6, != 1.5", "2.0", false},

		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L344-L356
		{">= 1.4.4, < 1.5", "1.4.5", true},
		{">= 1.4.4 <1.5", "1.5.0.rc1", true},
		{">= 1.4.4, < 1.5", "1.5.0", false},

		{">= 1.4.4, < 1.5.a", "1.4.5", true},
		{">= 1.4.4, < 1.5.a", "1.5.0.rc1", false},
		{">= 1.4.4  < 1.5.a", "1.5.0", false},

		// Bad
		// https://github.com/rubygems/rubygems/blob/v3.1.4/test/rubygems/test_gem_requirement.rb#L371-L386
		{"> 0.1", "", false},
		{"!= 1.2.3", "1.2.3", false},
		{"!= 1.02.3", "1.2.003.0.0", false},
		{"< 1.2.3", "4.5.6", false},
		{"> 1.1", "1.0", false},
		{"= 0.1", "", false},
		{"== 0.1", "", false},
		{"> 1.1.1", "1.1.1", false},
		{"= 1.1", "1.2", false},
		{"== 1.1", "1.2", false},
		{"= 1.1", "1.40", false},
		{"= 1.40", "1.3", false},
		{"<= 9.3.2", "9.3.3", false},
		{">= 9.3.2", "9.3.1", false},
		{"<= 9.3.2", "9.3.03", false},
		{"= 1.0", "1.0.0.1", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.version, tt.constraint), func(t *testing.T) {
			c, err := NewConstraints(tt.constraint)
			require.NoError(t, err)

			v, err := NewVersion(tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.want, c.Check(v))
		})
	}
}
