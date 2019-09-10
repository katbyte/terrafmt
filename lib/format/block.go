package format

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/katbyte/terrafmt/lib/common"
)

func Block(b string) (string, error) {

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

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

	return fb, nil
}
