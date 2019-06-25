package cli

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/katbyte/terrafmt/common"
)

func FormatBlock(b string, fmtCompat bool) (string, error) {

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// make block with fmt string placeholders tf fmt compatible
	if fmtCompat {

		// handle bare %s
		//figure this out later
		b = string(regexp.MustCompile(`(?m:^%s$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))
		b = string(regexp.MustCompile(`(?m:^[\s]*%s$)`).ReplaceAllString(b, `#@@_@@ TFMT:$0`))

		// handle [%s]

		b = string(regexp.MustCompile(`(?m:\[%s\])`).ReplaceAllString(b, `["@@_@@ TFMT:$0:TFMT @@_@@"]`))
	}

	common.Log.Debugf("running terraform... ")
	cmd := exec.Command("terraform", "fmt", "-")
	cmd.Stdin = strings.NewReader(b)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed with %s: %s", err, stderr)
	}

	ec := cmd.ProcessState.ExitCode()
	common.Log.Debugf("terraform exited with %d", ec)
	if ec != 0 {
		return "", fmt.Errorf("trerraform failed with %d: %s", ec, stderr)
	}
	fb := stdout.String()

	if fmtCompat {
		fb = strings.ReplaceAll(fb, "#@@_@@ TFMT:", "")
		fb = strings.ReplaceAll(fb, "[\"@@_@@ TFMT:", "")
		fb = strings.ReplaceAll(fb, ":TFMT @@_@@\"]", "")
	}

	return fb, nil
}
