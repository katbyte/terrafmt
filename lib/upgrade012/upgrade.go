package upgrade012

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/plugin"
	"github.com/katbyte/terrafmt/lib/common"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

func Block(b string) (string, error) {
	// Make temp directory
	dir, err := ioutil.TempDir(".", "tmp-module")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	// Create temp file
	tmpFile, err := ioutil.TempFile(dir, "*.tf")
	if err != nil {
		return "", err
	}

	// Write from Reader to File
	if _, err := tmpFile.Write(bytes.NewBufferString(b).Bytes()); err != nil {
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	// Empty the internal provisioners so as to avoid any setup job,
	// which is not necessary for upgrade job.
	command.InternalProvisioners = map[string]plugin.ProvisionerFunc{}

	// Buffer the upgrade ui output, to be part of the error message (if any)
	uiBuffer := bytes.NewBufferString("")
	meta := command.Meta{
		Ui: &cli.BasicUi{
			Writer: uiBuffer,
			Reader: nil,
		},
		GlobalPluginDirs: func() []string {
			ret := []string{}
			dir := filepath.Join(os.Getenv("HOME"), ".terraform.d")
			machineDir := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
			ret = append(ret, filepath.Join(dir, "plugins"))
			ret = append(ret, filepath.Join(dir, "plugins", machineDir))
			return ret
		}(),
	}

	upgradeCmd := &command.ZeroTwelveUpgradeCommand{Meta: meta}

	stdout, stderr, err := common.CaptureRun(func() error {
		// We specify "-force" here to allow a re-upgrade, otherwise it will complain the
		// cfg has already been upgraded.
		// Though, a re-upgrade will not always work, since it might failed to parse the 0.12 syntax
		// using 0.11 syntax parser.
		rc := upgradeCmd.Run([]string{"-yes", "-force", dir})
		if rc != 0 {
			return fmt.Errorf("upgrade to 0.12 failed (rc: %d): %s", rc, uiBuffer.String())
		}
		return nil
	})
	if viper.GetBool("verbose") {
		fmt.Fprintln(os.Stdout, stdout)
		fmt.Fprintln(os.Stderr, stderr)
	}
	if err != nil {
		return "", err
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
