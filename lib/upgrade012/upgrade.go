package upgrade012

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/katbyte/terrafmt/lib/common"
)

func Block(b string) (string, error) {

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Make temp directory
	dir, err := ioutil.TempDir(".", "tmp-module")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	common.Log.Debugf("running terraform... ")
	cmd := exec.Command("terraform", "0.12upgrade", "-yes", dir)

	// Create temp file
	tmpFile, err := ioutil.TempFile(dir, "*.tf")
	if err != nil {
		return "", err
	}

	defer os.Remove(tmpFile.Name()) // clean up

	// Write from Reader to File
	if _, err := tmpFile.Write(bytes.NewBufferString(b).Bytes()); err != nil {
		tmpFile.Close()
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed with %s: %s", err, stderr)
	}

	ec := cmd.ProcessState.ExitCode()
	common.Log.Debugf("terraform exited with %d", ec)
	if ec != 0 {
		return "", fmt.Errorf("terraform failed with %d: %s", ec, stderr)
	}

	// Read from temp file
	fb, err := ioutil.ReadFile(tmpFile.Name())

	return string(fb), nil
}
