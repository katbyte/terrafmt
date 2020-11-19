package upgrade012

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
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

	ctx := context.Background()

	tfBin, err := tfinstall.Find(ctx, tfinstall.ExactVersion("0.12.29", tempDir))
	if err != nil {
		return "", err
	}

	relTfBin, err := filepath.Rel(tempDir, tfBin)
	if err != nil {
		return "", err
	}

	tf, err := tfexec.NewTerraform(tempDir, relTfBin)
	if err != nil {
		return "", err
	}
	tf.SetStdout(stdout)
	tf.SetStderr(stderr)

	err = tf.Init(ctx)
	if err != nil {
		return "", fmt.Errorf("terraform init failed: %w", err)
	}

	err = tf.Upgrade012(ctx)
	if err != nil {
		return "", fmt.Errorf("terraform 0.12upgrade failed: %w", err)
	}

	// Read from temp file
	raw, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", tmpFile.Name(), err)
	}

	// 0.12upgrade always adds a trailing newline, even if it's already there
	// strip it here
	fb := string(raw)
	if strings.HasSuffix(fb, "\n\n") {
		fb = strings.TrimSuffix(fb, "\n")
	}

	return fb, nil
}
