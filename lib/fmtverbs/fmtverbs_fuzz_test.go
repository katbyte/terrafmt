package fmtverbs

import (
	"strings"
	"testing"
)

// FuzzEscapeUnscapeRoundTrip checks the core property of this package: escaping a block
// and unescaping it again must return the original input, byte for byte.
func FuzzEscapeUnscapeRoundTrip(f *testing.F) {
	seeds := []string{
		"",
		"%s\n",
		"    %d\n",
		"\t%q\n",
		"%[1]s\n",
		"%.1f\n",
		"attr = %s\n",
		"attr = %[7]q\n",
		"attr = [%s, %s]\n",
		"attr = [%[1]s, %[2]s]\n",
		"count = %d\n",
		"provider = %s\n",
		"provider = %[3]s\n",
		"cond = %[1]t ? 1 : 0\n",
		"attr = \"${map.%s.prop}\"\n",
		"attr = \"${list[%d]}\"\n",
		"attr = something.%s.prop\n",
		"attr = something.text%[2]s.prop\n",
		"attr = function(%s, %d)\n",
		"resource \"azurerm_thing\" \"%s\" {\n  name = %q\n}\n",
		"data \"azurerm_thing\" \"test-%[1]s\" {\n}\n",
		"list \"azurerm_thing\" \"%s\" {\n}\n",
		"ephemeral \"azurerm_key_vault_secret\" \"%s\" {\n}\n",
		"action \"azurerm_virtual_machine_run_command\" %q {\n}\n",
		"resource \"resource\" \"test\" {\n  kat = \"byte\"\n}\n",
		"%v %x %.10f\n", // verbs the escaper does not handle
		"100%s\n",
		"%%s\n",
		"%s",       // no trailing newline
		"%s\r\n",   // CRLF
		"a\x00b\n", // NUL byte

		// regressions found by fuzzing; keep so every plain `go test` run replays them
		"provider = %[0]s0\n", // literal digit after an indexed provider verb blurred the marker boundary
		"[%s",                 // literal [ before a line-ending verb: the list-form unescape ate the [
		"%.0f0",               // precision verb with trailing text: Ω marker dropped/never restored the precision
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, block string) {
		// inputs already containing escape-marker text cannot round-trip by design;
		// real terraform blocks never contain these sentinels
		for _, marker := range []string{"TFMT", "TMFT", "tfmt", "Ω", "Ω", "@@_@@"} {
			if strings.Contains(block, marker) {
				t.Skip()
			}
		}

		escaped := Escape(block)
		roundtrip := Unscape(escaped)
		if roundtrip != block {
			t.Errorf("did not roundtrip:\n  input:     %q\n  escaped:   %q\n  roundtrip: %q", block, escaped, roundtrip)
		}
	})
}
