// Copyright (c) 2012-2016 The go-diff authors. All rights reserved.
// https://github.com/sergi/go-diff
// See the included LICENSE file for license details.
//
// go-diff is a Go implementation of Google's Diff, Match, and Patch library
// Original library is Copyright (c) 2006 Google Inc.
// http://code.google.com/p/google-diff-match-patch/

// Package diffmatchpatch offers robust algorithms to perform the operations required for synchronising plain text.
//
// Absorbed from github.com/sergi/go-diff (via github.com/katbyte/sergi-go-diff,
// MIT licensed, see LICENSE), trimmed to the diff path used by terrafmt; the
// match and patch operations have been removed.
package diffmatchpatch

import (
	"time"
)

// DiffMatchPatch holds the configuration for diff-match-patch operations.
type DiffMatchPatch struct {
	// Number of seconds to map a diff before giving up (0 for infinity).
	DiffTimeout time.Duration
	// Cost of an empty edit operation in terms of edit characters.
	DiffEditCost int
}

// New creates a new DiffMatchPatch object with default parameters.
func New() *DiffMatchPatch {
	// Defaults.
	return &DiffMatchPatch{
		DiffTimeout:  time.Second,
		DiffEditCost: 4,
	}
}
