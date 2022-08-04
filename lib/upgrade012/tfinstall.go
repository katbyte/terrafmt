package upgrade012

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-version"
	hcinstall "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
)

func InstallTerraform(ctx context.Context) (string, error) {
	vc, _ := version.NewConstraint("=0.12.29")
	rv := &releases.Versions{
		Product:     product.Terraform,
		Constraints: vc,
	}

	versions, err := rv.List(ctx)
	if err != nil {
		return "", err
	}

	installer := hcinstall.NewInstaller()

	if len(versions) != 1 {
		return "", fmt.Errorf("expected a single version but got %d", len(versions))
	}

	execPath, err := installer.Ensure(ctx, []src.Source{
		versions[0],
	})
	if err != nil {
		return "", err
	}

	absBin, err := filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	return absBin, nil
}
