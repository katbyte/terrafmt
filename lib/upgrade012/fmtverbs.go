package upgrade012

import (
	"github.com/katbyte/terrafmt/lib/fmtverbs"
	"github.com/sirupsen/logrus"
)

func Upgrade12VerbBlock(log *logrus.Logger, b string) (string, error) {
	b = fmtverbs.Escape(b)

	fb, err := Block(log, b)
	if err != nil {
		return fb, err
	}

	fb = fmtverbs.Unscape(fb)

	return fb, nil
}
