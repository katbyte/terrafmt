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
		tmpFile.Close()
		os.RemoveAll(tempDir)
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		os.RemoveAll(tempDir)
		log.Fatal(err)
	}

	defer os.RemoveAll(tempDir)

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
	// fb := strings.TrimSuffix(string(raw), "\n")

	return fb, nil
}
