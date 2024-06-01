package testcase

import (
	"testing"
)

func RunTestCases[TestcaseT any](t *testing.T, testcases map[string]TestcaseT, verify func(t *testing.T, testcase *TestcaseT)) {

	t.Parallel()

	for name, testcase := range testcases {
		testcase := &testcase
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			verify(t, testcase)
		})
	}
}
