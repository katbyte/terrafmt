package cli

import (
	"regexp"
	"testing"
)

var logMsgRegexp *regexp.Regexp

func init() {
	logMsgRegexp = regexp.MustCompile(".*msg=\"([^\"]+)\".*\n")
}

func checkExpectedErrors(t *testing.T, casename, errOutput string, expectedErrs []string) {
	if expectedErrCount := len(expectedErrs); expectedErrCount > 0 {
		allMatches := logMsgRegexp.FindAllStringSubmatch(errOutput, -1)
		actualErrs := make([]string, 0, len(allMatches))
		for _, matches := range allMatches {
			if len(matches) == 2 {
				actualErrs = append(actualErrs, matches[1])
			}
		}
		if expectedErrCount != len(actualErrs) {
			t.Errorf("Case %q: Expected %d error messages\n%#v,\ngot %d\n%#v", casename, expectedErrCount, expectedErrs, len(actualErrs), actualErrs)
		} else {
			for i := range actualErrs {
				match, err := regexp.MatchString(regexp.QuoteMeta(expectedErrs[i]), actualErrs[i])
				if err != nil {
					t.Fatalf("Case %q, error message %d: error parsing regexp: %s", casename, i+1, err)
				}
				if !match {
					t.Errorf("Case %q: error message %d does not have match,\nexpected %q,\ngot      %q", casename, i+1, expectedErrs[i], actualErrs[i])
				}
			}
		}
	} else {
		if errOutput != "" {
			t.Errorf("Case %q: Got unexpected error output:\n%s", casename, errOutput)
		}
	}
}
