package diff_test

import (
	"testing"

	"github.com/katbyte/terrafmt/lib/diff"
)

// The expected outputs are pinned from github.com/katbyte/andreyvit-diff v0.0.3,
// which this package absorbed; byte-identical output was verified against that
// library by differential fuzzing and a corpus replay before it was removed.
func TestLineDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{
			name: "changed line",
			a:    "a\nb\nc\n",
			b:    "a\nB\nc\n",
			want: " a\n-b\n+B\n c",
		},
		{
			name: "removed line",
			a:    "a\nb\nc\n",
			b:    "a\nc\n",
			want: " a\n-b\n c",
		},
		{
			name: "added line",
			a:    "a\nc\n",
			b:    "a\nb\nc\n",
			want: " a\n+b\n c",
		},
		{
			name: "terraform block formatting",
			a:    "resource \"aws_thing\" \"test\" {\nname=\"foo\"\ncount=1\n}\n",
			b:    "resource \"aws_thing\" \"test\" {\n  name  = \"foo\"\n  count = 1\n}\n",
			want: " resource \"aws_thing\" \"test\" {\n-name=\"foo\"\n-count=1\n+  name  = \"foo\"\n+  count = 1\n }",
		},
		{
			name: "empty to content",
			a:    "",
			b:    "a\n",
			want: "+a",
		},
		{
			name: "no trailing newline",
			a:    "a",
			b:    "b",
			want: "-a\n+b",
		},
		{
			name: "multibyte runes",
			a:    "héllo wörld\n",
			b:    "héllo world\n",
			want: "-héllo wörld\n+héllo world",
		},
		{
			name: "identical input",
			a:    "a\nb\n",
			b:    "a\nb\n",
			want: " a\n b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := diff.LineDiff(tt.a, tt.b); got != tt.want {
				t.Errorf("LineDiff(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
