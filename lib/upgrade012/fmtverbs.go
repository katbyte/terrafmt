package upgrade012

import (
	"context"

	"github.com/katbyte/terrafmt/lib/fmtverbs"
	"github.com/sirupsen/logrus"
)

func Upgrade12VerbBlock(ctx context.Context, tfPath string, log *logrus.Logger, b string) (string, error) {
	b = fmtverbs.Escape(b)

	fb, err := Block(ctx, tfPath, log, b)
	if err != nil {
		return fb, err
	}

	fb = fmtverbs.Unscape(fb)

	return fb, nil
}
