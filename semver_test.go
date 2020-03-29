package semver

import (
	"reflect"
	"testing"
)

func TestVersions_Sort(t *testing.T) {
	tt := []struct {
		name     string
		sort     Versions
		expected Versions
	}{
		{
			name: "default",
			sort: Versions{
				{1, 2, 4},
				{2, 3, 4},
				{1, 2, 3},
				{1, 3, 4},
			},
			expected: Versions{
				{1, 2, 3},
				{1, 2, 4},
				{1, 3, 4},
				{2, 3, 4},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.sort.Sort(); !reflect.DeepEqual(tc.sort, tc.expected) {
				t.Errorf("wrong order: %v", tc.sort)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	refs := map[int]string{
		lt: "less",
		eq: "equal",
		gt: "greater",
	}

	tt := []struct {
		name       string
		comparator Version
		comparable Version
		expected   int
	}{
		{
			name:       "equal",
			comparator: Version{1, 2, 3},
			comparable: Version{1, 2, 3},
			expected:   eq,
		},
		// patch
		{
			name:       "greater: major: eq, minor: eq; patch: gt",
			comparator: Version{1, 2, 4},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "less: major: eq, minor: eq; patch: lt",
			comparator: Version{1, 2, 3},
			comparable: Version{1, 2, 4},
			expected:   lt,
		},
		// minor
		{
			name:       "greater: major: eq, minor: gt; patch: eq",
			comparator: Version{1, 3, 3},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "greater: major: eq, minor: gt; patch: lt",
			comparator: Version{1, 3, 3},
			comparable: Version{1, 2, 4},
			expected:   gt,
		},
		{
			name:       "less: major: eq, minor: lt; patch: eq",
			comparator: Version{1, 2, 3},
			comparable: Version{1, 3, 3},
			expected:   lt,
		},
		{
			name:       "less: major: eq, minor: lt; patch: gt",
			comparator: Version{1, 2, 4},
			comparable: Version{1, 3, 3},
			expected:   lt,
		},
		// major
		{
			name:       "greater: major: gt, minor: eq; patch: eq",
			comparator: Version{2, 2, 3},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "greater: major: gt, minor: gt; patch: eq",
			comparator: Version{2, 3, 3},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "greater: major: gt, minor: gt; patch: gt",
			comparator: Version{2, 3, 4},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "greater: major: gt, minor: lt; patch: gt",
			comparator: Version{2, 1, 4},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "greater: major: gt, minor: lt; patch: lt",
			comparator: Version{2, 1, 2},
			comparable: Version{1, 2, 3},
			expected:   gt,
		},
		{
			name:       "less: major: lt, minor: eq; patch: eq",
			comparator: Version{1, 2, 3},
			comparable: Version{2, 2, 3},
			expected:   lt,
		},
		{
			name:       "less: major: lt, minor: gt; patch: eq",
			comparator: Version{1, 3, 3},
			comparable: Version{2, 2, 3},
			expected:   lt,
		},
		{
			name:       "less: major: gt, minor: gt; patch: gt",
			comparator: Version{1, 3, 4},
			comparable: Version{2, 2, 3},
			expected:   lt,
		},
		{
			name:       "less: major: gt, minor: lt; patch: gt",
			comparator: Version{1, 1, 4},
			comparable: Version{2, 2, 3},
			expected:   lt,
		},
		{
			name:       "less: major: gt, minor: lt; patch: lt",
			comparator: Version{1, 1, 2},
			comparable: Version{2, 2, 3},
			expected:   lt,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tc.comparator.Compare(tc.comparable); actual != tc.expected {
				t.Errorf("unexpected comparison result: %s", refs[actual])
			}
		})
	}
}

func Test_Parse(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		expected Version
		err      error
	}{
		{
			name:     "empty string",
			s:        "",
			expected: Version{},
			err:      ErrEmpty,
		},
		{
			name:     "invalid format: excess element",
			s:        "comparator.2.3.4",
			expected: Version{},
			err:      ErrInvalidFormat,
		},
		{
			name:     "invalid format: not enough elements",
			s:        "v1.2",
			expected: Version{},
			err:      ErrInvalidFormat,
		},
		{
			name:     "empty element",
			s:        "v1.2.",
			expected: Version{},
			err:      ErrEmptyElement,
		},
		{
			name:     "invalid format: invalid character",
			s:        "v1.2.x",
			expected: Version{},
			err:      ErrInvalidCharacter,
		},
		{
			name:     "ok",
			s:        "v1.2.3",
			expected: Version{1, 2, 3},
			err:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Parse(tc.s)
			if err != tc.err {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expected.Eq(actual) {
				t.Errorf("expected: %s, but got: %s", tc.expected, actual)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tt := []struct {
		name     string
		v        Version
		expected string
	}{
		{
			name: "default",
			v: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: "v1.2.3",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.v.String()
			if actual != tc.expected {
				t.Errorf("expected: %s, but got: %s", tc.expected, actual)
			}
		})
	}
}
