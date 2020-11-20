package upgrade012

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfinstall"
)

func InstallTerraform(ctx context.Context) (string, error) {
	bin, err := tfinstall.Find(ctx, tfinstall.ExactVersion("0.12.29", ""))
	if err != nil {
		return "", err
	}

	absBin, err := filepath.Abs(bin)
	if err != nil {
		return "", err
	}
	return absBin, nil
}
