package semver

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var (
	// ErrEmpty means that given value is empty.
	ErrEmpty = errors.New("semver: empty string")
	// ErrEmptyElement means that some elements of version is empty.
	ErrEmptyElement = errors.New("semver: empty element")
	// ErrInvalidFormat means that given value to parse isn't in semver format.
	ErrInvalidFormat = errors.New("semver: invalid format")
	// ErrLeadingZeroes means that some elements of version has leading zeroes.
	ErrLeadingZeroes = errors.New("semver: leading zeros")
	// ErrInvalidCharacter means that there is some unsupported characters in version specification.
	ErrInvalidCharacter = errors.New("semver: invalid character")
)

const (
	major int = iota
	minor
	patch
)

const (
	numbers = "0123456789"
)

var zero = Version{}

type Versions []Version

func (vs Versions) Len() int           { return len(vs) }
func (vs Versions) Less(i, j int) bool { return vs[i].Lt(vs[j]) }
func (vs Versions) Swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs Versions) Sort()              { sort.Sort(vs) }
func (vs Versions) Ascending()         { vs.Sort() }
func Sort(vs ...Version) Versions      { sorted := Versions(vs); sorted.Sort(); return sorted }

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v Version) String() string {
	b := make([]byte, 0, 6)
	b = append(b, 'v')
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Patch, 10)
	return string(b)
}

const (
	lt = iota - 1
	eq
	gt
)

func (v Version) Eq(to Version) bool  { return v.Compare(to) == eq }
func (v Version) Lt(to Version) bool  { return v.Compare(to) == lt }
func (v Version) Lte(to Version) bool { return v.Compare(to) <= eq }
func (v Version) Gt(to Version) bool  { return v.Compare(to) == gt }
func (v Version) Gte(to Version) bool { return v.Compare(to) >= eq }
func (v Version) Compare(to Version) int {
	switch {
	case v.Major > to.Major:
		return gt
	case v.Major < to.Major:
		return lt
	case v.Minor > to.Minor:
		return gt
	case v.Minor < to.Minor:
		return lt
	case v.Patch > to.Patch:
		return gt
	case v.Patch < to.Patch:
		return lt
	default:
		return eq
	}
}

func New(s string) (Version, error) { return Parse(s) }

func V(s string) Version { return Must(s) }
func Must(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

func Parse(s string) (v Version, err error) {
	if len(s) == 0 {
		return zero, ErrEmpty
	}
	parts := strings.Split(strings.TrimPrefix(s, "v"), ".")
	if len(parts) != 3 {
		return zero, ErrInvalidFormat
	}
	if empty(parts...) {
		return zero, ErrEmptyElement
	}
	if !containsOnly(numbers, parts...) {
		return zero, ErrInvalidCharacter
	}
	if leadingZeroes(parts...) {
		return zero, ErrLeadingZeroes
	}
	if v.Major, err = strconv.ParseUint(parts[major], 10, 64); err != nil {
		return zero, err
	}
	if v.Minor, err = strconv.ParseUint(parts[minor], 10, 64); err != nil {
		return zero, err
	}
	if v.Patch, err = strconv.ParseUint(parts[patch], 10, 64); err != nil {
		return zero, err
	}
	return v, nil
}

func empty(vs ...string) bool {
	for _, v := range vs {
		if v == "" {
			return true
		}
	}
	return false
}

func leadingZeroes(vs ...string) bool {
	for _, v := range vs {
		if len(v) > 1 && v[0] == '0' {
			return true
		}
	}
	return false
}

func containsOnly(set string, vs ...string) bool {
	for _, v := range vs {
		if strings.IndexFunc(v, func(r rune) bool {
			return !strings.ContainsRune(set, r)
		}) > -1 {
			return false
		}
	}
	return true
}
