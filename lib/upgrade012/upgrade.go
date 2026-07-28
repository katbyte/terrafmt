// Package upgrade012 upgrades terraform blocks to 0.12 syntax by running them
// through a real terraform binary's 0.12upgrade command.
package upgrade012

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
)

func Block(ctx context.Context, tfPath string, log *logrus.Logger, b string) (string, error) {
	// Make temp directory
	tempDir, err := os.MkdirTemp(".", "tmp-module")
	if err != nil {
		log.Fatal(err)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp(tempDir, "*.tf")
	if err != nil {
		return "", err
	}

	// Write from Reader to File
	if _, err := tmpFile.Write(bytes.NewBufferString(b).Bytes()); err != nil {
		if closeErr := tmpFile.Close(); closeErr != nil {
			log.Warnf("failed to close temp file %s: %v", tmpFile.Name(), closeErr)
		}
		removeTempDir(log, tempDir)
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		removeTempDir(log, tempDir)
		log.Fatal(err)
	}

	defer removeTempDir(log, tempDir)

	tf, err := tfexec.NewTerraform(tempDir, tfPath)
	if err != nil {
		return "", err
	}

	err = tf.Init(ctx)
	if err != nil {
		return "", fmt.Errorf("terraform init failed: %w", err)
	}

	err = tf.Upgrade012(ctx)
	if err != nil {
		return "", fmt.Errorf("terraform 0.12upgrade failed: %w", err)
	}

	// Read from temp file
	raw, err := os.ReadFile(tmpFile.Name())
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

func removeTempDir(log *logrus.Logger, tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		log.Warnf("failed to remove temp dir %s: %v", tempDir, err)
	}
}
