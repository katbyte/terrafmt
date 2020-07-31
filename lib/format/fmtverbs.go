package format

import (
	"github.com/katbyte/terrafmt/lib/fmtverbs"
	"github.com/sirupsen/logrus"
)

func FmtVerbBlock(log *logrus.Logger, content, path string) (string, error) {
	content = fmtverbs.Escape(content)

	fb, err := Block(log, content, path)
	if err != nil {
		return fb, err
	}

	fb = fmtverbs.Unscape(fb)

	return fb, nil
}
