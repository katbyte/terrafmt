package cli

import (
	"regexp"
	"strings"
	"testing"
)

func checkExpectedErrors(t *testing.T, casename, errOutput string, expectedErrs []string) {
	if expectedErrCount := len(expectedErrs); expectedErrCount > 0 {
		actualErrs := strings.FieldsFunc(errOutput, func(c rune) bool {
			return c == '\n'
		})
		if expectedErrCount != len(actualErrs) {
			t.Errorf("Case %q: Expected %d error messages\n%#v,\ngot %d\n%#v", casename, expectedErrCount, expectedErrs, len(actualErrs), actualErrs)
		} else {
			for i := range actualErrs {
				match, err := regexp.MatchString(regexp.QuoteMeta(expectedErrs[i]), actualErrs[i])
				if err != nil {
					t.Fatalf("Case %q, error message %d: error parsing regexp: %s", casename, i+1, err)
				}
				if !match {
					t.Errorf("Case %q: error message %d no match,\nexpected %q,\ngot      %q", casename, i+1, expectedErrs[i], actualErrs[i])
				}
			}
		}
	} else {
		if errOutput != "" {
			t.Errorf("Case %q: Got unexpected error output:\n%s", casename, errOutput)
		}
	}
}
