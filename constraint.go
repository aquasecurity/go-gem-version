package gem

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

var (
	constraintOperators = map[string]operatorFunc{
		"":   constraintEqual,
		"=":  constraintEqual,
		"!=": constraintNotEqual,
		">":  constraintGreaterThan,
		"<":  constraintLessThan,
		">=": constraintGreaterThanEqual,
		"=>": constraintGreaterThanEqual,
		"<=": constraintLessThanEqual,
		"=<": constraintLessThanEqual,
		"~>": constraintPessimistic,
	}
	constraintRegexp *regexp.Regexp
)

type operatorFunc func(v, c Version) bool

func init() {
	ops := make([]string, 0, len(constraintOperators))
	for k := range constraintOperators {
		ops = append(ops, regexp.QuoteMeta(k))
	}

	constraintRegexp = regexp.MustCompile(fmt.Sprintf(
		`^\s*(%s)\s*(%s)\s*$`,
		strings.Join(ops, "|"),
		versionPattern))
}

// Constraints is one or more constraint that a version can be checked against.
type Constraints [][]constraint

type constraint struct {
	version  Version
	operator operatorFunc
}

// NewConstraints parses a given constraint and returns a new instance of Constraints
func NewConstraints(v string) (Constraints, error) {
	var css [][]constraint
	for _, vv := range strings.Split(v, "||") {
		var cs []constraint
		for _, single := range strings.Split(vv, ",") {
			c, err := newConstraint(single)
			if err != nil {
				return Constraints{}, err
			}
			cs = append(cs, c)
		}
		css = append(css, cs)
	}

	return css, nil

}

func newConstraint(c string) (constraint, error) {
	m := constraintRegexp.FindStringSubmatch(c)
	if m == nil {
		return constraint{}, xerrors.Errorf("improper constraint: %s", c)
	}

	v, err := NewVersion(m[2])
	if err != nil {
		return constraint{}, xerrors.Errorf("version parse error (%s): %w", m[2], err)
	}

	return constraint{
		version:  v,
		operator: constraintOperators[m[1]],
	}, nil
}

func (c constraint) check(v Version) bool {
	return c.operator(v, c.version)
}

// Check tests if a version satisfies all the constraints.
func (cs Constraints) Check(v Version) bool {
	for _, c := range cs {
		if andCheck(v, c) {
			return true
		}
	}

	return false
}

func andCheck(v Version, constraints []constraint) bool {
	for _, c := range constraints {
		if !c.check(v) {
			return false
		}
	}
	return true
}

//-------------------------------------------------------------------
// Constraint functions
//-------------------------------------------------------------------

func constraintEqual(v, c Version) bool {
	return v.Equal(c)
}

func constraintNotEqual(v, c Version) bool {
	return !v.Equal(c)
}

func constraintGreaterThan(v, c Version) bool {
	return v.GreaterThan(c)
}

func constraintLessThan(v, c Version) bool {
	return v.LessThan(c)
}

func constraintGreaterThanEqual(v, c Version) bool {
	return v.GreaterThanOrEqual(c)
}

func constraintLessThanEqual(v, c Version) bool {
	return v.LessThanOrEqual(c)
}

func constraintPessimistic(v, c Version) bool {
	return v.GreaterThanOrEqual(c) && v.Release().LessThan(c.Bump())
}