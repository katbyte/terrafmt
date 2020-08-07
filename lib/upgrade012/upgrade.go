package upgrade012

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func Block(log *logrus.Logger, b string) (string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Make temp directory
	tempDir, err := ioutil.TempDir(".", "tmp-module")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempDir) // clean up

	// Create temp file
	tmpFile, err := ioutil.TempFile(tempDir, "*.tf")
	if err != nil {
		return "", err
	}

	// Write from Reader to File
	if _, err := tmpFile.Write(bytes.NewBufferString(b).Bytes()); err != nil {
		tmpFile.Close()
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("terraform", "init")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(),
		"TF_IN_AUTOMATION=1",
	)
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed in terraform init with %s: %s", err, stderr)
	}

	log.Debugf("running terraform... ")
	cmd = exec.Command("terraform", "0.12upgrade", "-yes")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(),
		"TF_IN_AUTOMATION=1",
	)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()

	if err != nil {
		_, err := fmt.Println(stdout)
		if err != nil {
			return "", fmt.Errorf("cmd.Run() failed in terraform 0.12upgrade with %s: %s | %s", err, stdout, stderr)
		}

		return "", fmt.Errorf("cmd.Run() failed in terraform 0.12upgrade with %s: %s", err, stderr)
	}

	ec := cmd.ProcessState.ExitCode()
	log.Debugf("terraform exited with %d", ec)
	if ec != 0 {
		return "", fmt.Errorf("terraform failed with %d: %s", ec, stderr)
	}

	// Read from temp file
	raw, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("terrafmt failed with readfile: %s", err)
	}

	// 0.12upgrade always adds a trailing newline, even if it's already there
	// strip it here
	fb := string(raw)
	if strings.HasSuffix(fb, "\n") {
		fb = strings.TrimSuffix(fb, "\n")
	}

	return fb, nil
}
