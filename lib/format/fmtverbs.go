package format

import (
	"github.com/katbyte/terrafmt/lib/fmtverbs"
)

func FmtVerbBlock(b string) (string, error) {
	b = fmtverbs.Escape(b)

	fb, err := Block(b)
	if err != nil {
		return fb, err
	}

	fb = fmtverbs.Unscape(fb)

	return fb, nil
}
