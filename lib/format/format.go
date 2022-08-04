package format

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
)

func Block(log *logrus.Logger, content, path string) (string, error) {
	b := []byte(content)
	log.Debugf("format terraform config... ")
	_, syntaxDiags := hclsyntax.ParseConfig(b, path, hcl.Pos{Line: 1, Column: 1})
	if syntaxDiags.HasErrors() {
		return "", fmt.Errorf("failed to parse hcl: %w\n%s", errors.New(syntaxDiags.Error()), b)
	}

	return string(hclwrite.Format(b)), nil
}
